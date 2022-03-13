import Card from "react-bootstrap/Card";
import {useEffect, useState} from "react";
import {FormLabel, FormControl, FormText, FormSelect, FormGroup, Button} from "react-bootstrap";
import {timezones} from "./timezones.js";

function safePayload(timeZone, schedule, command) {
    return encodeURIComponent(JSON.stringify({
        timeZone, schedule, command
    }));
}

export default function Generator() {
    var options = timezones.map(tz => <option value={tz} key={tz}>{tz}</option>);
    var [timeZone, setTimeZone] = useState("America/New_York");
    var [schedule, setSchedule] = useState("");
    var [command, setCommand] = useState("");
    var [payload, setPayload] = useState("");

    useEffect(() => {
        setPayload(safePayload(timeZone, schedule, command));
    }, [timeZone, schedule, command]);

    const submitCallback = (e) => {
        e.preventDefault();
        alert(`${timeZone} ${schedule} ${command}`);
    };

    return <>
        <Card className={"generator-section p-3"}>
            <form onSubmit={submitCallback}>
                <h2>Options</h2>

                <FormGroup className="my-3">
                    <FormLabel htmlFor="timezoneInput">Time Zone</FormLabel>
                    <FormSelect id="timezoneInput" value={timeZone} onChange={(e) => setTimeZone(e.target.value)}>
                        {options}
                    </FormSelect>
                </FormGroup>

                <FormGroup className="my-3">
                    <FormLabel htmlFor="scheduleInput">Schedule</FormLabel>
                    <FormControl
                        type="text"
                        id="scheduleInput"
                        aria-describedby="scheduleInputHelpBlock"
                        placeholder="e.g. */5 * * * * for every 5 minutes."
                        value={schedule}
                        onChange={(e)=>setSchedule(e.target.value)}
                    />
                    <FormText id="scheduleInputHelpBlock">
                        Enter a cron compatible schedule. There is a helper tool available at&nbsp;
                        <a href="http://crontab.cronhub.io" target="_blank">crontab.cronhub.io.</a>
                    </FormText>
                </FormGroup>

                <FormGroup className="my-3">
                    <FormLabel htmlFor="commandInput">Command to Run</FormLabel>
                    <FormControl
                        type="text"
                        id="commandInput"
                        aria-describedby="commandInputHelpBlock"
                        placeholder="e.g. curl -o /status.txt http://example.com/status"
                        value={command}
                        onChange={(e) => setCommand(e.target.value)}
                    />
                    <FormText id="commandInputHelpBlock">
                        The command as you would enter it into your terminal emulator.
                    </FormText>
                </FormGroup>

                <div className="right-align">
                    <a href={"/generate?payload=" + payload} target="_blank" className={"btn btn-primary"} download>Generate</a>
                </div>
            </form>
        </Card>
    </>;
}