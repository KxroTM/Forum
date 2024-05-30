function reportToggle(id) {
    var report = document.getElementById(id);
    console.log(id);
    if (report.style.opacity == 0) {
        report.style.opacity = 1;
        report.style.zIndex = 1;
    } else {
        report.style.opacity = 0;
        report.style.zIndex = -10;
    }
}