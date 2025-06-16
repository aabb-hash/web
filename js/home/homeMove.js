const circle = document.getElementById("user");
const keysDown = {};
let x = 50;
let y = 50;
const speed = 1;

const maxX = window.innerWidth - 20
const maxY = window.innerHeight - 20

function updatePosition(element, newX, newY) {
    if (element == null) {
        if (newX != null && newY != null) {
            x = newX
            y = newY
        }

        circle.style.left = x + "px"
        circle.style.top = y + "px"

        sendCoords(x, y)
        return
    }

    element.style.left = newX + "px";
    element.style.top = newY + "px";
}

function wait(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}
function noBarrier(dirX, dirY) {
    const offsets = [
        [20, 20],
        [-20, 20],
        [20, -20],
        [-20, -20]
    ];
    const finalX = x + dirX;
    const finalY = y + dirY;

    if (finalX < 20 || finalX > maxX || finalY < 20 || finalY > maxY) {
        return false
    }

    return offsets.every(([dx, dy]) => {
        const elements = document.elementsFromPoint(finalX + dx, finalY + dy);
        return !elements.some(el => el.classList.contains("barrier"));
    });
}
async function keepMoving(dir) {
    switch (dir) {
        case "up":
            while (keysDown["up"]) {
                if (!keysDown["down"] && noBarrier(0, -1)) {
                    y -= speed;
                    updatePosition();
                }
                await wait(1);
            }
            break;
        case "down":
            while (keysDown["down"]) {
                if (!keysDown["up"] && noBarrier(0, 1)) {
                    y += speed;
                    updatePosition();
                }
                await wait(1);
            }
            break;
        case "left":
            while (keysDown["left"]) {
                if (!keysDown["right"] && noBarrier(-1, 0)) {
                    x -= speed;
                    updatePosition();
                }
                await wait(1);
            }
            break;
        case "right":
            while (keysDown["right"]) {
                if (!keysDown["left"] && noBarrier(1, 0)) {
                    x += speed;
                    updatePosition();
                }
                await wait(1);
            }
            break;
    }
}

function getKey(key) {
    switch (key) {
        case "ArrowUp":
        case "w":
        case "W":
            return "up"
        case "ArrowDown":
        case "s":
        case "S":
            return "down"
        case "ArrowLeft":
        case "a":
        case "A":
            return "left"
        case "ArrowRight":
        case "d":
        case "D":
            return "right"
    }
}

window.addEventListener("keydown", (e) => {
    if (e.repeat) return;

    key = getKey(e.key)
    if (!keysDown[key]) {
        keysDown[key] = true
        keepMoving(key)
    }
});

window.addEventListener('keyup', (e) => {
    key = getKey(e.key)
    keysDown[key] = false;
});