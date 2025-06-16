fetch("/api/leaderboard")
  .then(response => {
    if (!response.ok) {
      throw new Error("Failed to load leaderboard data.");
    }
    return response.json();
  })
  .then(players => {
    const tbody = document.querySelector("#leaderboard tbody");
    players.forEach(player => {
      const row = document.createElement("tr");
      row.innerHTML = `
        <td>${player.name}</td>
        <td>${player.wins}</td>
        <td>${player.losses}</td>
        <td>${player.draws}</td>
      `;
      tbody.appendChild(row);
    });
  })
  .catch(error => {
    console.error("Error loading leaderboard:", error);
  });