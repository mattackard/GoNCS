let metricsContainer = document.getElementById("metricsContainer");
let serviceContainers = document.getElementsByClassName("serviceContainer");
let serviceTitles = document.getElementsByClassName("serviceTitle");
let serviceInfo = document.getElementsByClassName("serviceInfo");

let logContainer = document.getElementById("logContainer");
let logFilterContainer = document.getElementById("logFilterContainer");
let logFilters = document.getElementsByClassName("logFilter");
let logDisplay = document.getElementById("logDisplay");

//get the title of the log file (todays date)
let date = new Date(Date.now());
let day = date.getDate() < 10 ? `0${date.getDate()}` : date.getDate();
let month = date.getMonth() + 1 < 10 ? `0${date.getMonth() + 1}` : date.getMonth() + 1;
let logTitle = `${date.getFullYear()}-${month}-${day}.txt`;

//tell the go server to get the master log and record it into the file
setInterval(() => {
    fetch("http://localhost:80/getLogs").then(response => response.text()).then(text => {
        console.log("log response " + text);
        logDisplay.value = text;
    });
}, 5000);

//tell the go server to get the master log and record it into the file
setInterval(() => {
    fetch("http://localhost:80/stats").then(response => response.text()).then(json => {
        console.log("stat response " +json);
    });
}, 5000);