        function togglePostType() {
            var postType = document.querySelector('input[name="post_type"]:checked').value;
            var textPostFields = document.getElementById('text-post-fields');
            var imagePostFields = document.getElementById('image-post-fields');
            if (postType === 'text') {
                textPostFields.style.display = 'block';
                imagePostFields.style.display = 'none';
            } else {
                textPostFields.style.display = 'none';
                imagePostFields.style.display = 'block';
            }
        }
        window.onload = function() {
            togglePostType(); // Set initial state based on default selection
        };