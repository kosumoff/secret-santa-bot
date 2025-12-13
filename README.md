# Secret Santa Bot 🎅

A lightweight Telegram bot for organizing gift exchanges, built with Go and Clean Architecture.

## ✨ Features

* **Private Registration**: Users join via deep-linking to keep the group chat clean.
* **Smart Draw**: Uses a circular shift algorithm to ensure everyone gives and receives exactly one gift.
* **Admin Control**: Only group administrators can trigger the shuffle.
* **Direct Notifications**: Results are sent privately to each participant.

## 💻 Tech Stack
* **Language**: Go 1.21+
* **Framework**: [go-telegram/bot](https://github.com/go-telegram/bot)
* **Database**: SQLite
* **Architecture**: Clean Architecture (Domain, Usecase, Adapter layers)

## 🚀 Quick Start

1. **Clone & Setup**:
   ```bash
   git clone https://github.com/kosumoff/secret-santa-bot.git
   cd secret-santa-bot
   ```

2. **Configure `.env`**:
   ```env
   BOT_TOKEN=your_bot_api_token_here
   DB_PATH=./santa.db
   ```

3. **Run**:
   ```bash
   go run cmd/bot/main.go
   ```

## 🛠 Commands

* `/santa` — Post the "Join" button in a group.
* `/draw` — (Admin only) Shuffle and send assignments.
* `/start` — Register the user via private message.

## 🎮 How to Use

1. **Add the bot** to your Telegram group.
2. **Promote as Admin**: Ensure it has permissions to send messages.
3. **Start Joining**: Type `/santa` in the group to post the registration button.
4. **Register**: Click the "Join Secret Santa" button. You will be redirected to a private chat with the bot. Press the "Start" button there to confirm your participation.
5. **Draw**: Once everyone has joined, an admin types `/draw` in the group. The bot will send the names of recipients to everyone via private message.

## 📄 License

This project is licensed under the [MIT License](LICENSE).