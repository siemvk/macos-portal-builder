from enum import Enum
import tkinter as tk


class LogLevel(Enum):
    DEBUG = 3
    WARNING = 2
    INFO = 1
    ERROR = 0


class Logger:
    def __init__(self, log_file, log_level: LogLevel = LogLevel.INFO):
        self.log_file = log_file
        self.log_level = log_level if isinstance(log_level, LogLevel) else LogLevel.INFO

    def debug(self, message):
        if self.log_level.value >= LogLevel.DEBUG.value:
            self.log(f"\033[37m[DEBUG] {message}\033[0m")

    def openTkinterLogger(self):
        self.tkinter_window = tk.Tk()
        self.tkinter_window.title("Logger")
        self.tkinter_window.geometry("400x300")
        self.tkinter_text = tk.Text(self.tkinter_window)
        self.tkinter_text.pack(expand=True, fill=tk.BOTH)

    def warn(self, message):
        if self.log_level.value >= LogLevel.WARNING.value:
            self.log(f"\033[33m[WARNING]\033[0m {message}")

    def warning(self, message):
        self.warn(message)

    def info(self, message):
        if self.log_level.value >= LogLevel.INFO.value:
            self.log(f"\033[34m[INFO]\033[0m {message}")

    def error(self, message):
        if self.log_level.value >= LogLevel.ERROR.value:
            self.log(f"\033[31m[ERROR]\033[0m {message}")

    def success(self, message):
        if self.log_level.value >= LogLevel.INFO.value:
            self.log(f"\033[32m[SUCCESS]\033[0m {message}")

    def log(self, message):
        print(str(message))
        # with open(self.log_file, "a", encoding="utf-8") as f:
        #     f.write(str(message) + "\n")
        if hasattr(self, "tkinter_text"):
            self.tkinter_text.insert(tk.END, str(message) + "\n")
            self.tkinter_text.see(tk.END)


class logging(Logger):
    pass
