
setInterval(function() {
    location.reload();
}, 10000); 

let moveHistory = [];  // Track the history of moves
let currentFEN = document.getElementById("board").getAttribute("game-fen");
console.log("Current FEN:", currentFEN); // Debugging: Check if FEN is being retrieved correctly


// Function to process the FEN string and create the board
function processFEN(fen) {
    const fenRows = fen.split(' ')[0].split('/');
    const board = [];

    for (let row of fenRows) {
        const rowArray = [];
        let col = 0;

        for (let char of row) {
            if (parseInt(char)) {
                // If the character is a number (empty squares), add that many empty spaces
                for (let i = 0; i < parseInt(char); i++) {
                    rowArray.push(' ');  // Add an empty space for each empty square
                }
                col += parseInt(char);  // Skip the number of columns based on the digit
            } else {
                // Otherwise, it's a piece, so add it to the row and increment col
                rowArray.push(char);
                col++;
            }
        }

        board.push(rowArray);  // Add the row to the board array
    }
    // For debugging, logs the board structure
    return board;
}

// Function to generate the chessboard from the FEN notation
function createBoard(fen) {
    const boardContainer = document.getElementById('board');
    console.log(fen)
    const board = processFEN(fen);  // Process FEN into a board array
    // Clear any existing buttons
    boardContainer.innerHTML = '';

    // Loop through the board rows and columns
    for (let row = 0; row < 8; row++) {
        for (let col = 0; col < 8; col++) {
            const button = document.createElement('button');
            button.classList.add('board-button');

            // Alternate the colors for each square
            if ((row + col) % 2 === 0) {
                button.classList.add('white');
            } else {
                button.classList.add('black');
            }

            // Get the piece for the current square from the board array
            const piece = board[row][col];

            // If there's a piece on this square, set the background image accordingly
            if (piece !== undefined) {
                const pieceColor = piece === piece.toLowerCase() ? 'b' : 'w';
                const imageUrl = `figures/standard/${pieceColor}${piece.toLowerCase()}.png`;
                button.style.backgroundImage = `url(${imageUrl})`;
            }

            // Set the position and handle click
            button.dataset.position = `${String.fromCharCode(97 + col)}${8 - row}`;

            button.onclick = function() {
                handleMove(button);
            };

            // Add the button to the board container
            boardContainer.appendChild(button);
        }
    }
}

// Handle the move logic
function handleMove(button) {
    // Add the selected position to the move history
    if (moveHistory.length === 0 || moveHistory.length === 1) {
        moveHistory.push(button.dataset.position);

        // Update the "Moves" section with the selected square
        document.getElementById("moves").innerText = moveHistory.join(' to ');

        // If two moves are selected, send them to the server
        if (moveHistory.length === 2) {
            sendMoveToServer(moveHistory[0], moveHistory[1]);
            moveHistory = []; // Clear after sending
        }
    }
}

// Function to update the chessboard based on the FEN string
function updateBoard(fen) {
    const boardContainer = document.getElementById('board');
    const board = processFEN(fen);  // Process the FEN into a board array

    // Loop through the board rows and columns
    for (let row = 0; row < 8; row++) {
        for (let col = 0; col < 8; col++) {
            const button = boardContainer.children[row * 8 + col];  // Get the corresponding button

            // Reset the background image before setting a new one
            button.style.backgroundImage = '';

            // Get the piece for the current square from the board array
            const piece = board[row][col];

            // If there's a piece on this square, set the background image accordingly
            if (piece !== ' ') {  // Check if it's not an empty square (' ' is now used for empty spaces)
                const pieceColor = piece === piece.toLowerCase() ? 'b' : 'w';
                const imageUrl = `figures/standard/${pieceColor}${piece.toLowerCase()}.png`;
                button.style.backgroundImage = `url(${imageUrl})`;
            }
        }
    }
}

function sendMoveToServer(from, to) {
    fetch('/game', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ from, to, fen: currentFEN })
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();  // Attempt to parse JSON
    })
    .then(data => {
        if (data) {
            if (data.valid) {
                currentFEN = data.fen;  // Assuming the server returns the updated FEN
                console.log("Move successful");

                // Update the board with the new FEN
                updateBoard(currentFEN);
            } else {
                showErrorMessage();
            }
        } else {
            console.error('Received empty or invalid data:', data);
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
}

// Show the "Illegal Move" message
function showErrorMessage() {
    const errorMessage = document.getElementById('errorMessage');
    errorMessage.style.display = 'block';
    setTimeout(() => {
        errorMessage.style.display = 'none';
    }, 1000);
}

// Fetch the FEN string passed from the Go server (injected into the HTML template)
createBoard(currentFEN);  // Generate the board using FEN
