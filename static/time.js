function fetchBeijingTime() {
    fetch('http://worldtimeapi.org/api/timezone/Asia/Shanghai')
        .then(response => response.json())
        .then(data => {
            const dateTime = new Date(data.datetime);
            document.getElementById('current-time').innerText = `Current Beijing Time: ${dateTime.toLocaleString()}`;
            document.getElementById('current-timestamp').innerText = `Current Timestamp: ${dateTime.getTime()}`;
        })
        .catch(error => console.error('Error fetching Beijing time:', error));
}

function convertTimestampToTime(timestamp) {
    const dateTime = new Date(parseInt(timestamp));
    return dateTime.toLocaleString();
}

function convertTimeToTimestamp(timeString) {
    const dateTime = new Date(timeString);
    return dateTime.getTime();
}

document.getElementById('timestamp-to-time-form').addEventListener('submit', function(event) {
    event.preventDefault();
    const timestamp = document.getElementById('timestamp-input').value;
    const convertedTime = convertTimestampToTime(timestamp);
    document.getElementById('converted-time').innerText = `Converted Time: ${convertedTime}`;
});

document.getElementById('time-to-timestamp-form').addEventListener('submit', function(event) {
    event.preventDefault();
    const timeString = document.getElementById('time-input').value;
    const convertedTimestamp = convertTimeToTimestamp(timeString);
    document.getElementById('converted-timestamp').innerText = `Converted Timestamp: ${convertedTimestamp}`;
});

fetchBeijingTime();
setInterval(fetchBeijingTime, 60000);
