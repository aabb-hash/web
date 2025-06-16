const openSettingsBtn = document.getElementById('openSettings');
const closeSettingsBtn = document.getElementById('closeModal');
const settingsModal = document.getElementById('settingsModal');
const settingsForm = document.getElementById('settingsForm');
let open = false;

openSettingsBtn.addEventListener('click', () => {
    settingsModal.style.display = 'flex';
    open = true
});

closeSettingsBtn.addEventListener('click', () => {
    settingsModal.style.display = 'none';
    open = false
});

window.addEventListener('keydown', (e) => {
    if (!open) return;

    if (e.key == 'Escape') {
        settingsModal.style.display = 'none';
        open = false
    }
});

const input = document.getElementById('newAvatar');
const preview = document.getElementById('preview');

input.addEventListener('change', function () {
  const file = this.files[0];
  if (!file) return;

  if (!['image/png', 'image/jpeg'].includes(file.type)) {
    alert('Only PNG and JPEG files are allowed.');
    this.value = '';
    preview.style.display = 'none';
    return;
  }

  const reader = new FileReader();
  reader.onload = function (e) {
    preview.src = e.target.result;
    preview.style.display = 'flex';
  };
  reader.readAsDataURL(file);
});


settingsForm.addEventListener("submit", function (e) {
  e.preventDefault()

  const formData = new FormData(settingsForm)

  fetch("/api/settings", {
    method: "PUT",
    body: formData
  })
  .then(async res => {
    if (res.status != 204) {
      const text = await res.text();
      throw new Error(`Server error: ${res.status} ${text}`);
    }
    window.location.replace(window.location)
    return res.text();
  })
  .catch(error => console.log("Request failed: ", error))
});
