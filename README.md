# Chess Web - Version 2.1 (DISCARDED PROJECT)

## Overview
Chess Web is a simple multiplayer chess game that allows you to play with a friend or by yourself. The game runs on a Go server and uses a web-based interface for making moves.

## Running the Game
1. Ensure you have Go installed.
2. Open a terminal and navigate to the project directory.
3. Run the following command:
   ```sh
   go run .
   ```
4. Open your browser and go to:
   ```
   http://localhost:8080
   ```
##THIS IS DISCARDED PROJECT (Second Warning)

## How to Play
1. Enter a unique 6-letter long game code (ONLY English letters are allowed). Case does not matter, as all letters will be converted to uppercase.
2. Once inside the game, you can play chess against a friend or by yourself.
3. To make a move:
   - Click on the first square (the piece you want to move).
   - Click on the second square (where you want to move the piece).
   - If you accidentally click the second square before the first, click again to cancel.
   - If the move is legal, it will be executed; otherwise, it will not be played.
4. If you experience lags, try reloading the page and making the move again. However, this may not always resolve the issue.

## Important Notes
- **Anyone who knows the game code can join and play, including taking over your opponent's moves.**
- **Pawn promotion does not work (but castling does).**
- **En passant is not verified to work.**
- **The following features are NOT implemented:**
  - Offering a draw, resigning, or scoring points.
  - Automatic draws (e.g., 50-move rule, stalemate, checkmate, time control).
  - Account system for players.
  - History of moves (PGN notation).

## Version History
- **1.0 - February 21, 21:50**
- **2.0 - March 1, 15:45**
- **2.1 - March 2, 17:15**
- **Project is discarded since it is too laggy**

---
Enjoy the game and forget this project

**Author:**
Ansar Zeinulla
