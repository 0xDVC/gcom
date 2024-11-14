# 🔍 gcom - Git Commit Observer & Manager
A simple command-line tool that helps you find and visualize commits across multiple Git repositories in your workspace. The idea came from missing a commit in a repo I was working on but forgot to push. I forgot both the name of the repo and the commit message. Does this happen to you? If it does, this was my solution. An overly engineered solution for a simple problem. An excuse to read golang docs and learn the language, with the good aid of gpt for illustrations and better understanding. Oh welp, it's a work-in-progress(WIP).

![Terminal Preview](https://via.placeholder.com/800x400?text=Terminal+Preview)

## ✨ Features

- 🎯 Find commits across multiple repositories in one go
- 📅 Filter commits by time range (today, hours, days)
- 🔄 Track pushed/unpushed commits
- 🎨 Beautiful terminal UI with color-coded output
- 📱 Responsive pager interface for easy navigation
- ⚡ Concurrent repository scanning for speed
- 📁 Recursively scan for Git repositories in subdirectories


## 🚀 Installation (DEVELOPMENT)

1. Clone the repository:
   ```sh
   git clone https://github.com/0xdvc/gcom.git
   ```

2. Navigate to the project directory:
   ```sh
   cd gcom
   ```

3. Build the project:
   ```sh
   go build -o gcom cmd/gcom/main.go
   ```

4. Run the tool:
   ```sh
   ./gcom
   ```


## 📝 License

MIT License - feel free to use and modify!