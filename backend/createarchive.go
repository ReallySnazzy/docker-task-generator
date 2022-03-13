package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

const crontabHeader = `
# ┌───────────── minute (0 - 59)
# │ ┌───────────── hour (0 - 23)
# │ │ ┌───────────── day of month (1 - 31)
# │ │ │ ┌───────────── month (1 - 12)
# │ │ │ │ ┌───────────── day of week (0 - 6) (Sunday to Saturday;
# │ │ │ │ │                                       7 is also Sunday on some systems)
# │ │ │ │ │
# │ │ │ │ │
# * * * * *  command to execute
`

func createCrontabFile(schedule string, command string) (string, error) {
	crontabFile, err := os.CreateTemp("", "crontab")
	if err != nil {
		return "", err
	}
	defer crontabFile.Close()
	crontabFile.WriteString(crontabHeader)
	crontabFile.WriteString(schedule)
	crontabFile.WriteString(" ")
	crontabFile.WriteString(command)
	crontabFile.WriteString("\n")
	return crontabFile.Name(), nil
}

func timeZoneDockerString(timeZone string) string {
	return fmt.Sprintf("%s%s%s%s",
		"RUN apk update && apk add tzdata && \\\n",
		fmt.Sprintf("\tcp /usr/share/zoneinfo/%s /etc/localtime && \\\n", timeZone),
		fmt.Sprintf("\techo \"%s\" > /etc/timezone && \\\n", timeZone),
		"\tapk del tzdata && rm -rf /var/cache/apk/*\n")
}

func createDockerfile(timeZone string) (string, error) {
	dockerFile, err := os.CreateTemp("", "Dockerfile")
	if err != nil {
		return "", err
	}
	defer dockerFile.Close()
	dockerFile.WriteString("# Put your previous build stages here\n\n")
	dockerFile.WriteString("FROM alpine:3.14\n\n")
	dockerFile.WriteString(timeZoneDockerString(timeZone))
	dockerFile.WriteString("\n")
	dockerFile.WriteString("COPY ./crontab /etc/crontabs/root\n\n")
	dockerFile.WriteString("# Copy your files from previous build stages here\n")
	dockerFile.WriteString("# COPY --from=0 /YOURFILE /\n\n")
	dockerFile.WriteString("CMD chown root:root /etc/crontabs/root && /usr/sbin/crond -f\n")
	return dockerFile.Name(), nil
}

func copyFileToZip(zipWriter *zip.Writer, sourceFileName string, zipFileName string) error {
	file, err := os.Open(sourceFileName)
	if err != nil {
		return err
	}
	defer file.Close()
	zipFileWriter, err := zipWriter.Create(zipFileName)
	if err != nil {
		return err
	}
	io.Copy(zipFileWriter, file)
	return nil
}

func addFilesToNewArchive(crontabFileName string, dockerFileName string) (string, error) {
	zipArchive, err := os.CreateTemp("", "package")
	if err != nil {
		return "", err
	}
	defer zipArchive.Close()

	writer := zip.NewWriter(zipArchive)
	defer writer.Close()

	err = copyFileToZip(writer, crontabFileName, "crontab")
	if err != nil {
		return "", err
	}

	err = copyFileToZip(writer, dockerFileName, "Dockerfile")
	if err != nil {
		return "", err
	}

	return zipArchive.Name(), nil
}

func createArchive(payload ArchiveRequestPayload) (string, error) {
	crontabFileName, err := createCrontabFile(payload.Schedule, payload.Command)
	if err != nil {
		return "", err
	}
	defer os.Remove(crontabFileName)

	dockerFileName, err := createDockerfile(payload.TimeZone)
	if err != nil {
		return "", err
	}
	defer os.Remove(dockerFileName)

	return addFilesToNewArchive(crontabFileName, dockerFileName)
}
