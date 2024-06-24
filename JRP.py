# Library imports
import customtkinter
from tkinter import messagebox
import subprocess


# Jason L use CTkinter finish that, feel free to edit and use it 😎.
# Project URL: https://github.com/JasonL111/Planner
# CTKinter: https://customtkinter.tomschimansky.com


# Save Textbox content to .env file
def create_popup(message):
    text_content = Textbox.get("1.0", customtkinter.END)
    filename = f".env"
    try:
        with open(filename, "w", encoding="utf-8") as file:
            file.write(text_content)
        messagebox.showinfo("Success", "File has been saved.")
    except Exception as e:
        messagebox.showerror("Error", f"Error while saving the file: {e}")


# Load .env from directory
def readfile():
    Textbox.delete("1.0", "end")
    filename = ".env"
    try:
        with open(filename, "r", encoding="utf-8") as file:
            content = file.read()
            Textbox.insert("1.0", content)
        messagebox.showinfo("Success", "File has been loaded.")
    except Exception as e:
        messagebox.showerror("Error", f"Error while loading the file: {e}")


# Delete current content in Textbox
def deleteplan():
    Textbox.delete("1.0", "end")
    Textbox.insert("1.0", "S3_KEY_ID=\n")
    Textbox.insert("2.0", "S3_SECRET=\n")
    Textbox.insert("3.0", "S3_REGOIN=\n")
    Textbox.insert("4.0", "S3_ENDPOINT=\n")
    Textbox.insert("5.0", "S3_BUCKET=\n")
    Textbox.insert("6.0", "S3_PREFIX=\n")
    Textbox.insert("7.0", "S3_DURATION_HOURS=\n")

# Execute main.exe and redirect the output to result.txt
# def generate():
#     subprocess.run(["go", "mod", "tidy"], check=True)
#     build_process = subprocess.run(["go", "build"], capture_output=True, text=True)
#     if build_process.returncode != 0:
#         messagebox.showerror("Error", f"Go build failed:\n{build_process.stderr}")
#         return
#     result = subprocess.run(['main.exe'], stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
#     with open('result.txt', 'w', encoding='utf-8') as file:
#         file.write(result.stdout)
#     if result.returncode == 0:
#         messagebox.showinfo("Success", "Command executed successfully and output saved to result.txt")
#     else:
#         messagebox.showerror("Error", f"Command failed with error: {result.stderr}")
def generate():
    try:
        subprocess.run(["final.bat"], check=True)
        messagebox.showinfo("Success", "Batch file executed successfully")
    except subprocess.CalledProcessError as e:
        messagebox.showerror("Error", f"Error while executing batch file: {e}")
    except FileNotFoundError:
        messagebox.showerror("Error", "Batch file not found")

# Initialize and configure main window
root = customtkinter.CTk()
root.geometry("600x650")
root.title("JR Presigned tool")
root.iconbitmap("icon.ico")
root.configure(bg="#FFFFF6")


# Create a label in the window
label = customtkinter.CTkLabel(root, text="JR Presigned Tool", font=("微软雅黑", 22))


# Create a Textbox for user input
Textbox = customtkinter.CTkTextbox(
    root, height=500, font=("微软雅黑", 18), width=500, corner_radius=10
)
Textbox.insert("1.0", "S3_KEY_ID=\n ")
Textbox.insert("2.0", "S3_SECRET=\n")
Textbox.insert("3.0", "S3_REGOIN=\n")
Textbox.insert("4.0", "S3_ENDPOINT=\n")
Textbox.insert("5.0", "S3_BUCKET=\n")
Textbox.insert("6.0", "S3_PREFIX=\n")
Textbox.insert("7.0", "S3_DURATION_HOURS=\n")
# Create buttons for functions
button1 = customtkinter.CTkButton(
    root,
    text="Save",
    font=("微软雅黑", 24),
    command=lambda: create_popup("Save Successful"),
    corner_radius=15,
    width=48,
    height=10,
)
button2 = customtkinter.CTkButton(
    root,
    text="Delete",
    font=("微软雅黑", 24),
    command=deleteplan,
    corner_radius=15,
    width=48,
    height=10,
)
button3 = customtkinter.CTkButton(
    root,
    text="Load",
    font=("微软雅黑", 24),
    command=readfile,
    corner_radius=15,
    width=48,
    height=10,
)
button4 = customtkinter.CTkButton(
    root,
    text="Generate",
    font=("微软雅黑", 24),
    command=generate,
    corner_radius=15,
    width=48,
    height=10,
)
for i in range(12):
    root.columnconfigure(i, weight=1)
for j in range(15):
    root.rowconfigure(j, weight=1)


# Place a label in the window
label.grid(column=0, row=0, columnspan=12)


# Place a Textbox for user input
Textbox.grid(column=1, row=1, columnspan=10, rowspan=11, sticky="nsew")

# Place buttons
button1.grid(column=1, row=13, columnspan=1, sticky="nsew")
button2.grid(column=4, row=13, columnspan=1, sticky="nsew")
button3.grid(column=7, row=13, columnspan=1, sticky="nsew")
button4.grid(column=10, row=13, columnspan=1, sticky="nsew")


# Start the main event loop of the window
root.mainloop()
