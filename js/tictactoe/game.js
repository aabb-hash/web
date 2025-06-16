const board = document.getElementById('board');
const gameID = window.location.pathname.split("/")[2]
const positionList = [
    "top-left",
    "top",
    "top-right",
    "mid-left",
    "mid",
    "mid-right",
    "bottom-left",
    "bottom",
    "bottom-right"
]

let opponentUsername
let opponentID
let turn = false

function loadMe() {
    const me = document.getElementById("me")

    const username = document.createElement("span")
    username.textContent = getCookie("username")

    me.appendChild(username)
}

function loadOpponent(newOpponentID, newOpponentUsername) {
    opponentID = newOpponentID
    opponentUsername = newOpponentUsername

    const opponent = document.getElementById("opponent")

    const avatar = document.createElement("img")
    avatar.src = "/api/avatar?id=" + newOpponentID
    avatar.alt = "Opponent"

    const username = document.createElement("span")
    username.textContent = newOpponentUsername

    opponent.appendChild(username)
    opponent.appendChild(avatar)

    loadStats(newOpponentUsername)
}

function loadBoard() {
    board.innerHTML = ""

    for (let i = 0; i < 9; i++) {
        const cell = document.createElement('div');
        cell.classList.add('cell');
        cell.classList.add(positionList.at(i))

        cell.addEventListener('click', () => {
            if (turn && !cell.classList.contains('x') && !cell.classList.contains('o')) {
                cell.classList.add('x');

                document.getElementById("me").classList.remove("active")
                document.getElementById("opponent").classList.add("active")

                turn = false
                sendMove(i)
            }
        });

        board.appendChild(cell);
    }
}

function loadStats(opponent) {
    fetch("/api/stats-to-opponent?opponent=" + opponent)
        .then(response => response.text())
        .then(data => {
            stats = data.split(";")
            document.getElementById("wins").textContent = data[0];
            document.getElementById("draws").textContent = data[2];
            document.getElementById("losses").textContent = data[4];
        })
}

function opponentMove(position) {
    const cell = document.getElementsByClassName(positionList.at(position))[0]
    cell.classList.add('o');

    document.getElementById("opponent").classList.remove("active")
    document.getElementById("me").classList.add("active")

    turn = true
}

function gameUpdate(signal, data) {
    switch (signal) {
        case 0:
            document.getElementById('resultPopup').classList.add('hidden');
            loadBoard(null)
            
            turn = false
            document.getElementById("me").classList.remove("active")
            document.getElementById("opponent").classList.add("active")
            break
        case 1:
            document.getElementById('resultPopup').classList.add('hidden');
            loadBoard(null)

            turn = true
            document.getElementById("opponent").classList.remove("active")
            document.getElementById("me").classList.add("active")
            break
        case 2:
            showResultPopup(1)
            break
        case 3:
            showResultPopup(2)
            break
        case 4:
            showResultPopup(0)
            break
        case 5:
            updateOpponentStatus(1)
            break
        case 6:
            updateOpponentStatus(2)
            break
    }
}

function getCookie(name) {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) {
    return parts.pop().split(';').shift();
  }
  return null;
}

loadMe()
loadBoard()
