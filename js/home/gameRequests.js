const requestsSection = document.getElementById('game-requests-section');
const toggleBtn = document.getElementById('toggle-requests-btn');
const requests = new Map()

toggleBtn.addEventListener('click', () => {
  requestsSection.classList.toggle('closed');
});


function addRequest(id, game, type = 'received') {
  const container = document.getElementById('requests');
  const card = document.createElement('div');
  card.className = 'request-card animate-in';
  card.id = `${type}-request-${id}`;

  if (type == "sent") {
    game = requests.get(id)
  } else {
    requests.set(id, game)
  }

  let gameDisplay
  switch (game) {
    case 1:
      gameDisplay = "TicTacToe"
  }

  card.innerHTML = `
    <div class="request-top">
      <img src="/api/avatar?id=${id}" alt="${id}" class="request-avatar">
      <div class="request-username">${users.get(id)}</div>
    </div>
    <div class="request-info"><i class="fas fa-gamepad"></i> ${gameDisplay}</div>
    <div class="request-buttons">
      ${
        type === 'received'
          ? `
            <button class="accept-btn" onclick="handleAccept('${id}')">Accept</button>
            <button class="reject-btn" onclick="handleReject('${id}')">Reject</button>
          `
          : `<button class="cancel-btn" onclick="handleCancel('${id}')">Cancel</button>`
      }
    </div>
  `;

  container.appendChild(card);
  setTimeout(() => card.classList.remove('animate-in'), 300);
}

function removeRequest(id, type = 'received') {
  const card = document.getElementById(`${type}-request-${id}`);
  if (card) {
    card.remove();
  }
}

function handleAccept(id) {
  id = Number(id);

  removeRequest(id, 'received');

  const buffer = new ArrayBuffer(3);
  const view = new DataView(buffer);

  view.setUint8(0, 4)
  view.setUint8(1, requests.get(id))
  view.setUint8(2, id);

  sendMessage(buffer);
}

function handleReject(id) {
  console.log(`Rejected request ${id}`);
  removeRequest(id, 'received');
  // send to server...
}

function handleCancel(id) {
  console.log(`Canceled request ${id}`);
  removeRequest(id, 'sent');
  // send to server...
}



const modal = document.getElementById('create-request-modal');
const closeBtn = document.getElementById('close-create-modal');
const playerList = document.getElementById('player-list');
const searchInput = document.getElementById('player-search');
const gameModeSelect = document.getElementById('game-mode');

document.getElementById('open-create-request').addEventListener('click', () => {
  modal.classList.remove('hidden');
  renderPlayerList('');
});

closeBtn.addEventListener('click', () => {
  modal.classList.add('hidden');
});

window.addEventListener('keydown', (e) => {
    if (e.key == 'Escape' && !modal.classList.contains('hidden')) {
        !modal.classList.add('hidden')
    }
});

searchInput.addEventListener('input', () => {
  renderPlayerList(searchInput.value);
});

function renderPlayerList(searchTerm) {
  playerList.innerHTML = '';

  users.forEach((username, id) => {
    if (!username.toLowerCase().includes(searchTerm.toLowerCase())) return

    const item = document.createElement('div');
    item.className = 'player-item';
 
    item.innerHTML = `
      <img src="/api/avatar?id=${id}" alt="${id}" class="player-avatar" />
      <span>${username}</span>
    `;

    item.addEventListener('click', () => {
      const selectedMode = gameModeSelect.value;
      sendGameRequest(id, selectedMode);
      modal.classList.add('hidden');
    });
    playerList.appendChild(item);
  });
}

function sendGameRequest(playerId, mode) {
  const buffer = new ArrayBuffer(3);
  const view = new DataView(buffer);

  view.setUint8(0, 4)
  switch (mode) {
    case "TicTacToe":
      view.setUint8(1, 1)
      break
  }
  view.setUint8(2, playerId);

  requests.set(playerId, mode)
  sendMessage(buffer);
}