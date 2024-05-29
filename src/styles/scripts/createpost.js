        function togglePostType() {
            const textPostFields = document.getElementById('text-post-fields');
            const imagePostFields = document.getElementById('image-post-fields');
            const textPostRadio = document.getElementById('text-post');

            if (textPostRadio.checked) {
                textPostFields.style.display = 'block';
                imagePostFields.style.display = 'none';
                document.getElementById('title').required = true;
                document.getElementById('text').required = true;
                document.getElementById('image-title').required = false;
                document.getElementById('imageschoose').required = false;
            } else {
                textPostFields.style.display = 'none';
                imagePostFields.style.display = 'block';
                document.getElementById('title').required = false;
                document.getElementById('text').required = false;
                document.getElementById('image-title').required = true;
                document.getElementById('imageschoose').required = true;
            }
        }

        document.addEventListener('DOMContentLoaded', (event) => {
            togglePostType();
        });

        function validateForm() {
            const textPostFieldsVisible = document.getElementById('text-post-fields').style.display !== 'none';
            const imagePostFieldsVisible = document.getElementById('image-post-fields').style.display !== 'none';

            if (textPostFieldsVisible) {
                document.getElementById('image-title').removeAttribute('required');
                document.getElementById('imageschoose').removeAttribute('required');
            } else if (imagePostFieldsVisible) {
                document.getElementById('title').removeAttribute('required');
                document.getElementById('text').removeAttribute('required');
            }

            return true;
        }