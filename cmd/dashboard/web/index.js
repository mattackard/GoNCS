let metricsContainer = document.getElementById("metricsContainer");
let serviceContainers = document.getElementsByClassName("serviceContainer");
let serviceUL = document.getElementById("serviceUL");
let serviceTitles = document.getElementsByClassName("serviceTitle");
let serviceInfo = document.getElementsByClassName("serviceInfo");

let logContainer = document.getElementById("logContainer");
let logFilterContainer = document.getElementById("logFilterContainer");
let logFilters = document.getElementsByClassName("logFilter");
let logDisplay = document.getElementById("logDisplay");

let requestInterval = 5000;

//get the title of the log file (todays date)
let date = new Date(Date.now());
let day = date.getDate() < 10 ? `0${date.getDate()}` : date.getDate();
let month = date.getMonth() + 1 < 10 ? `0${date.getMonth() + 1}` : date.getMonth() + 1;
let logTitle = `${date.getFullYear()}-${month}-${day}.txt`;

//tell the go server to get the master log and record it into the file
setInterval(() => {
    fetch("http://localhost:80/getLogs").then(response => response.text()).then(text => {
        logDisplay.value = text;
    });
}, requestInterval);

//tell the go server to get the master log and record it into the file
setInterval(() => {
    fetch("http://localhost:80/stats").then(response => response.json()).then(json => {
        let html = "";
        json.containers.forEach(service => {
            html += createStatListItem(service);
        });
        serviceUL.innerHTML = html;
    });
}, requestInterval);

function createStatListItem(serviceStats) {
    let html = `<li>
                    <h3 class="serviceTitle">${serviceStats.serviceName}</h3>
                    <h3 class="ipAddress">IP: ${serviceStats.ip}</h3>
                    <ul class="serviceInfo">
                        <li>CPU Usage ${serviceStats.cpuShare}</li>
                        <li>CPU User Time ${serviceStats.cpuUserTime}</li>
                        <li>CPU System Time ${serviceStats.cpuSysTime}</li>
                        <li>Available Memory ${serviceStats.availableMem}</li>
                        <li>Memory Use ${serviceStats.memUsage}</li>
                        <li>Open Files ${serviceStats.openFiles}</li>
                        <li>Thread Count ${serviceStats.threadCount}</li>
                    </ul>
                </li>`;
    return html;
}