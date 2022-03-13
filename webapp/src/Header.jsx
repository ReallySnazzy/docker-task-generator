export default function Header(props) {
    var {title} = props;
    return (<>
        <header className="jumbo">
            <h1>{title}</h1>
            <p>
                Generate Dockerfiles that perform tasks using your schedule and timezone information.
            </p>
        </header>
    </>);
}