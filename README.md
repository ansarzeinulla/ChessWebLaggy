# Chess Web - Version 2.1

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

## Future Plans
1) Implement a custom chess game class to support up to **20x20** boards and **2-8 players**.
2) Add **offering draw, resign, and points tracking**.
3) Implement **automatic draws (50-move rule, stalemate, checkmate detection)**.
4) Enable **pawn promotion**.
5) Store **move history using PGN notation**.
6) Ensure **uniqueness of 6-letter game IDs**.

---
Enjoy the game and stay tuned for updates!

