document.addEventListener('DOMContentLoaded', function() {
    const repostCheckbox = document.getElementById('repostCheckbox');
    const repostPopup = document.querySelector('.repost-popup');

    repostCheckbox.addEventListener('change', function() {
        if (repostCheckbox.checked) {
            repostPopup.style.display = 'block';
        } else {
            repostPopup.style.display = 'none';
        }
    });

    document.addEventListener('click', function(event) {
        if (!repostCheckbox.contains(event.target) && !repostPopup.contains(event.target)) {
            repostCheckbox.checked = false;
            repostPopup.style.display = 'none';
        }
    });
});
