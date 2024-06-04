function reportToggle(id) {
    var report = document.getElementById(id+"report");
    if (report.style.opacity == 0) {
        report.style.opacity = 1;
        report.style.zIndex = 1;
    } else {
        report.style.opacity = 0;
        report.style.zIndex = -10;
    }
}

function reportOverlay(id) {
    var report = document.getElementById("reportOverlay");
    var button = document.getElementById("reportLink");
    if (report.style.opacity == 0) {
        report.style.opacity = 1;
        report.style.zIndex = 10000;
        if (id !== "") {
            button.href = "/report?"+id;
        }
    } else {
        report.style.opacity = 0;
        report.style.zIndex = -10;
    }
}

function toggleNotif() {
    fleche = document.getElementById("piqueoverlay");
    notif = document.getElementById("notifoverlay");

    if (notif.style.opacity == 0) {
        notif.style.opacity = 1;
        notif.style.zIndex = 10000;
        fleche.style.opacity = 1;
        fleche.style.zIndex = 10001;
    } else {
        notif.style.opacity = 0;
        notif.style.zIndex = -10;
        fleche.style.opacity = 0;
        notif.style.zIndex = -10;
    }
}

function reglageToggle() {
    var reglage = document.getElementById("fullreglage");
    if (reglage.style.opacity == 0) {
        reglage.style.opacity = 1;
        reglage.style.zIndex = 1000;
    } else {
        reglage.style.opacity = 0;
        reglage.style.zIndex = -10;
    }
}

function toggleProfile() {
    overlay = document.getElementById("overlayprofile");

    if (overlay.style.opacity == 0) {
        overlay.style.opacity = 1;
        overlay.style.zIndex = 10000;
    } else {
        overlay.style.opacity = 0;
        overlay.style.zIndex = -1000;
    }
}

document.addEventListener("DOMContentLoaded", () => {
    const postImages = document.querySelectorAll('#postimg');
    const post = document.querySelectorAll('.post');
    const colorMode = document.getElementById('theme-switch');

    postImages.forEach(img => {
        img.addEventListener('click', () => {
            console.log('click');
            const overlay = document.getElementById('overlay');
            const overlayimg = document.getElementById('imgoverlay');
            if (overlay.style.display === 'none') {
                overlay.style.display = 'flex';
                overlay.style.width = '100%';
                overlay.style.height = '100%';
                overlay.style.position = 'fixed';
                overlay.style.top = '0';
                overlay.style.left = '0';
                overlayimg.src = img.src;
                overlayimg.style.height = '65%';
                overlay.style.zIndex = '10000';
            }
            overlay.addEventListener('click', () => {
                overlay.style.display = 'none';
            });
        });
    });

    post.forEach(p => {
        p.addEventListener('click', () => {
            if (!event.target.classList.contains('ignore-click') && !event.target.closest('.ignore-click')) {
                window.location.href = '/post/id=' + p.id;
            }});
    });

    colorMode.addEventListener('click' , () => {
        setTimeout(() => {
            window.location.href = '/changeColorMode';
        }, 500);
    });

});


function likePost(id) {    
    window.location.href = '/likePost?'+id;
}

function dislikePost(id) {
    window.location.href = '/dislikePost?'+id;
}

function retweetPost(id) {
    window.location.href = '/retweetPost?'+id;
}