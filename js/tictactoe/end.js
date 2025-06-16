let opponentStatus = 0

function showResultPopup(result) {
    let myText
    let opponentText
    
    switch (result) {
        case 0:
            myText = "It's a draw!"
            opponentText = "drew the match"
            break
        case 1:
            myText = "You won!"
            opponentText = "Lost the match"
            break
        case 2:
            myText = "You lost!"
            opponentText = "Won the match"
            break
    }

    document.getElementById('resultPopup').classList.remove('hidden');
    document.getElementById('resultTitle').innerText = myText;
    document.getElementById('opponentText').innerText = `${opponentUsername} - ${opponentText}`;

    const opponentImg = document.querySelector('.opponent-status img');
    opponentImg.src = "/api/avatar?id=" + opponentID;
}

function updateOpponentStatus(status) {
    opponentStatus = status
    let opponentText
    
    switch (status) {
        case 1:
            opponentText = "Want's to play again"
            break
        case 2:
            opponentText = "Left"
            break
    }

    document.getElementById('opponentText').innerText = `${opponentUsername} - ${opponentText}`;
}


document.getElementById('returnHomeBtn').addEventListener('click', () => {
    const buffer = new ArrayBuffer(1);
    new DataView(buffer).setUint8(0, 6)
    sendMessage(buffer)

    window.location.replace(window.location.origin);
});

document.getElementById('playAgainBtn').addEventListener('click', () => {
    if (opponentStatus == 2) return

    document.getElementById('resultTitle').innerText = "waiting..";

    const buffer = new ArrayBuffer(2);
    const view = new DataView(buffer)

    view.setUint8(0, 4)
    view.setUint8(1, gameID)

    sendMessage(buffer)
});
