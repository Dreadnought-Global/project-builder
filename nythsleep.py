#!/usr/bin/env python3
"""
NythSleep - A stylish Windows 11 power management CLI tool.
Shutdown, restart, sleep, or logout with optional timers.
"""

import os
import sys
import time
import ctypes
import subprocess
import urllib.request
import json
import argparse

# ── ANSI Color Codes ──────────────────────────────────────────────────────────
RESET   = "\033[0m"
BOLD    = "\033[1m"
DIM     = "\033[2m"
UNDERLINE = "\033[4m"

# Purple gradient shades
P1 = "\033[38;2;200;140;255m"  # light lavender
P2 = "\033[38;2;180;100;255m"  # soft purple
P3 = "\033[38;2;160;70;255m"   # medium purple
P4 = "\033[38;2;140;50;230m"   # deep purple
P5 = "\033[38;2;120;30;210m"   # darker purple
P6 = "\033[38;2;100;20;190m"   # darkest purple

# Accent colors
CYAN    = "\033[38;2;100;220;255m"
PINK    = "\033[38;2;255;120;200m"
GREEN   = "\033[38;2;120;255;160m"
YELLOW  = "\033[38;2;255;230;100m"
RED     = "\033[38;2;255;100;100m"
WHITE   = "\033[38;2;220;220;230m"
GRAY    = "\033[38;2;130;130;150m"

def set_theme(name):
    global P1, P2, P3, P4, P5, P6
    if name == "midnight":
        P1, P2, P3, P4, P5, P6 = (
            "\033[38;2;140;200;255m", "\033[38;2;100;180;255m", "\033[38;2;70;160;255m",
            "\033[38;2;50;140;230m", "\033[38;2;30;120;210m", "\033[38;2;20;100;190m"
        )
    elif name == "sunset":
        P1, P2, P3, P4, P5, P6 = (
            "\033[38;2;255;180;140m", "\033[38;2;255;150;100m", "\033[38;2;255;120;70m",
            "\033[38;2;230;90;50m", "\033[38;2;210;60;30m", "\033[38;2;190;40;20m"
        )
    elif name == "forest":
        P1, P2, P3, P4, P5, P6 = (
            "\033[38;2;140;255;180m", "\033[38;2;100;255;150m", "\033[38;2;70;255;120m",
            "\033[38;2;50;230;90m", "\033[38;2;30;210;60m", "\033[38;2;20;190;40m"
        )
    else: # lavender
        P1, P2, P3, P4, P5, P6 = (
            "\033[38;2;200;140;255m", "\033[38;2;180;100;255m", "\033[38;2;160;70;255m",
            "\033[38;2;140;50;230m", "\033[38;2;120;30;210m", "\033[38;2;100;20;190m"
        )

VERSION = "1.2.0"
AUTHOR  = "NythSleep"

# ── ASCII Banner ──────────────────────────────────────────────────────────────

def get_banner():
    # ANSI Shadow Font: "NYTH SLEEP"
    return f"""
{P1}{BOLD}  ███╗   ██╗██╗   ██╗████████╗██╗  ██╗     ███████╗██╗     ███████╗███████╗██████╗ 
{P2}  ████╗  ██║╚██╗ ██╔╝╚══██╔══╝██║  ██║     ██╔════╝██║     ██╔════╝██╔════╝██╔══██╗
{P3}  ██╔██╗ ██║ ╚████╔╝    ██║   ███████║     ███████╗██║     █████╗  █████╗  ██████╔╝
{P4}  ██║╚██╗██║  ╚██╔╝     ██║   ██╔══██║     ╚════██║██║     ██╔══╝  ██╔══╝  ██╔═══╝ 
{P5}  ██║ ╚████║   ██║      ██║   ██║  ██║     ███████║███████╗███████╗███████╗██║     
{P6}  ╚═╝  ╚═══╝   ╚═╝      ╚═╝   ╚═╝  ╚═╝     ╚══════╝╚══════╝╚══════╝╚══════╝╚═╝     {RESET}
"""

def clear_screen():
    os.system("cls" if os.name == "nt" else "clear")

def print_banner(latest_version=None):
    print(get_banner())
    print(f"  {GRAY}{'─' * 104}{RESET}")
    link = f"\033]8;;https://fromrha.com\033\\{UNDERLINE}fromrha.com{RESET}\033]8;;\033\\"
    print(f"  {P2}{BOLD}{AUTHOR}{RESET} {GRAY}v{VERSION}{RESET}  {GRAY}|{RESET}  {WHITE}Windows Power Management Tool{RESET}  {GRAY}|{RESET}  {GRAY}{time.strftime('%Y-%m-%d')}{RESET}  {GRAY}|{RESET}  {P1}{link}{RESET}")
    print(f"  {GRAY}{'─' * 104}{RESET}")
    if latest_version:
        print(f"  {YELLOW}{BOLD}Update Available!{RESET} {WHITE}v{latest_version} is out. Download it from GitHub.{RESET}")
        print(f"  {GRAY}{'─' * 104}{RESET}")
    print()

def check_for_update():
    try:
        url = "https://api.github.com/repos/fromrha/nythsleep/releases/latest"
        req = urllib.request.Request(url, headers={'User-Agent': 'Nythsleep-Update-Checker'})
        with urllib.request.urlopen(req, timeout=1.0) as response:
            data = json.loads(response.read().decode())
            latest_version = data.get("tag_name", VERSION).lstrip('v')
            if latest_version != VERSION:
                return latest_version
    except Exception:
        pass
    return None

# ── Menu ──────────────────────────────────────────────────────────────────────

ACTIONS = [
    ("Shutdown",  "Power off your machine completely"),
    ("Restart",   "Reboot your machine"),
    ("Sleep",     "Put your machine to sleep"),
    ("Logout",    "Sign out of the current user session"),
]

def show_menu():
    print(f"  {P2}{BOLD}What would you like to do?{RESET}\n")
    for i, (name, desc) in enumerate(ACTIONS, 1):
        color = [CYAN, GREEN, YELLOW, PINK][i - 1]
        print(f"    {color}{BOLD}[{i}]{RESET}  {WHITE}{name}{RESET}  {GRAY}- {desc}{RESET}")
    print()
    print(f"    {RED}{BOLD}[0]{RESET}  {WHITE}Exit{RESET}")
    print()

def get_choice():
    while True:
        try:
            raw = input(f"  {P3}{BOLD}>{RESET} Pick an option {GRAY}(0-4){RESET}: ").strip()
            choice = int(raw)
            if 0 <= choice <= 4:
                return choice
            print(f"  {RED}  Please enter a number between 0 and 4.{RESET}")
        except (ValueError, EOFError):
            print(f"  {RED}  Invalid input. Enter a number.{RESET}")

# ── Timer ─────────────────────────────────────────────────────────────────────

def parse_timer(raw: str) -> int:
    """Parse timer string like '1h 30m 10s' into total seconds."""
    raw = raw.strip().lower()
    if raw in ("", "0", "now"):
        return 0

    total = 0
    current_num = ""
    for ch in raw:
        if ch.isdigit():
            current_num += ch
        elif ch in ("h", "m", "s") and current_num:
            val = int(current_num)
            if ch == "h":
                total += val * 3600
            elif ch == "m":
                total += val * 60
            elif ch == "s":
                total += val
            current_num = ""
        elif ch == " ":
            continue
        else:
            return -1  # invalid
    if current_num:
        return -1  # trailing digits with no unit
    return total

def format_duration(seconds: int) -> str:
    if seconds == 0:
        return f"{YELLOW}immediately{RESET}"
    parts = []
    h = seconds // 3600
    m = (seconds % 3600) // 60
    s = seconds % 60
    if h: parts.append(f"{h}h")
    if m: parts.append(f"{m}m")
    if s: parts.append(f"{s}s")
    return f"{YELLOW}{' '.join(parts)}{RESET}"

def get_timer():
    print(f"\n  {P3}{BOLD}>{RESET} Set a timer {GRAY}(e.g. 1h 30m, 45m, 10s, or press Enter for now){RESET}:")
    while True:
        raw = input(f"  {P3}{BOLD}>{RESET} ").strip()
        if raw == "":
            return 0
        secs = parse_timer(raw)
        if secs < 0:
            print(f"  {RED}  Invalid format. Use combinations like: 1h  30m  90s  1h 30m 10s{RESET}")
        else:
            return secs

# ── Pro Features ──────────────────────────────────────────────────────────────

ES_CONTINUOUS = 0x80000000
ES_SYSTEM_REQUIRED = 0x00000001
ES_DISPLAY_REQUIRED = 0x00000002

def set_insomnia_mode(enable: bool):
    if os.name == "nt":
        if enable:
            ctypes.windll.kernel32.SetThreadExecutionState(ES_CONTINUOUS | ES_SYSTEM_REQUIRED | ES_DISPLAY_REQUIRED)
        else:
            ctypes.windll.kernel32.SetThreadExecutionState(ES_CONTINUOUS)

def run_insomnia_loop(latest_version=None, duration_secs=0):
    clear_screen()
    print_banner(latest_version)
    
    status_text = f"{YELLOW}{BOLD}Active Indefinitely{RESET}"
    if duration_secs > 0:
        status_text = f"{YELLOW}{BOLD}Active for {format_duration(duration_secs)}{RESET}"

    print(f"  {status_text}")
    print(f"  {GRAY}{'─' * 104}{RESET}")
    print(f"  {WHITE}System and screen are now being kept awake.{RESET}")
    print(f"  {GRAY}Press {RESET}{RED}{BOLD}Ctrl+C{RESET}{GRAY} to stop and exit.{RESET}\n")
    
    start_time = time.time()
    try:
        while True:
            elapsed = int(time.time() - start_time)
            
            if duration_secs > 0:
                remaining = duration_secs - elapsed
                if remaining <= 0:
                    print(f"\n\n  {GREEN}Insomnia timer finished. System can now sleep.{RESET}")
                    break
                h, m, s = remaining // 3600, (remaining % 3600) // 60, remaining % 60
                label = "Time Remaining:"
                time_str = f"{RED}{h:02d}:{m:02d}:{s:02d}{RESET}"
            else:
                h, m, s = elapsed // 3600, (elapsed % 3600) // 60, elapsed % 60
                label = "Elapsed Time:"
                time_str = f"{CYAN}{h:02d}:{m:02d}:{s:02d}{RESET}"
            
            # Pulsing effect
            dot = f"{P2}●{RESET}" if elapsed % 2 == 0 else f"{GRAY}○{RESET}"
            
            line = f"\r  {dot}  {WHITE}{label}{RESET} {time_str}  "
            sys.stdout.write(line)
            sys.stdout.flush()
            time.sleep(1)
            
        set_insomnia_mode(False)
        print(f"\n  {P2}Goodbye!{RESET}\n")
        sys.exit(0)
    except KeyboardInterrupt:
        set_insomnia_mode(False)
        print(f"\n\n  {P2}Insomnia mode disabled. Goodbye!{RESET}\n")
        sys.exit(0)

class SYSTEM_POWER_STATUS(ctypes.Structure):
    _fields_ = [
        ("ACLineStatus", ctypes.c_byte),
        ("BatteryFlag", ctypes.c_byte),
        ("BatteryLifePercent", ctypes.c_byte),
        ("SystemStatusFlag", ctypes.c_byte),
        ("BatteryLifeTime", ctypes.c_ulong),
        ("BatteryFullLifeTime", ctypes.c_ulong)
    ]

def get_battery_percentage() -> int:
    if os.name == "nt":
        status = SYSTEM_POWER_STATUS()
        if ctypes.windll.kernel32.GetSystemPowerStatus(ctypes.byref(status)):
            return status.BatteryLifePercent
    return -1

def wait_for_battery(target_percentage: int, action_name: str):
    print(f"\n  {P2}{BOLD}Waiting for battery to reach {target_percentage}%...{RESET}\n")
    try:
        while True:
            current = get_battery_percentage()
            if current != -1 and current <= target_percentage:
                print(f"\n  {YELLOW}Battery reached {current}%. Proceeding with {action_name}.{RESET}")
                break
            
            line = f"\r  {CYAN}Current Battery: {current}%{RESET}  {GRAY}(Target: {target_percentage}%){RESET}  {GRAY}{action_name} pending...{RESET}  "
            sys.stdout.write(line)
            sys.stdout.flush()
            time.sleep(10)
        sys.stdout.write("\r" + " " * 80 + "\r")
        sys.stdout.flush()
    except KeyboardInterrupt:
        print(f"\n\n  {RED}{BOLD}Cancelled!{RESET} Action aborted by user.\n")
        sys.exit(0)

def send_notification(title: str, message: str):
    if os.name == "nt":
        ps_script = f'''
[void] [System.Reflection.Assembly]::LoadWithPartialName("System.Windows.Forms")
$objNotifyIcon = New-Object System.Windows.Forms.NotifyIcon
$objNotifyIcon.Icon = [System.Drawing.SystemIcons]::Information
$objNotifyIcon.BalloonTipTitle = "{title}"
$objNotifyIcon.BalloonTipText = "{message}"
$objNotifyIcon.Visible = $True
$objNotifyIcon.ShowBalloonTip(10000)
Start-Sleep -Seconds 3
'''
        subprocess.Popen(["powershell", "-NoProfile", "-WindowStyle", "Hidden", "-Command", ps_script], 
                         creationflags=subprocess.CREATE_NO_WINDOW)

# ── Countdown ─────────────────────────────────────────────────────────────────

def countdown(seconds: int, action_name: str):
    if seconds == 0:
        return
    print(f"\n  {P2}{BOLD}Countdown{RESET} {GRAY}(press Ctrl+C to cancel){RESET}\n")
    notified = False
    try:
        remaining = seconds
        while remaining > 0:
            if remaining == 60 and not notified:
                send_notification("NythSleep", f"{action_name} will execute in 60 seconds.")
                notified = True


            h = remaining // 3600
            m = (remaining % 3600) // 60
            s = remaining % 60
            bar_total = 30
            filled = int((1 - remaining / seconds) * bar_total)
            bar = f"{P4}{'█' * filled}{GRAY}{'░' * (bar_total - filled)}{RESET}"
            line = f"\r  {bar}  {CYAN}{h:02d}:{m:02d}:{s:02d}{RESET}  {GRAY}{action_name} pending...{RESET}  "
            sys.stdout.write(line)
            sys.stdout.flush()
            time.sleep(1)
            remaining -= 1
        sys.stdout.write("\r" + " " * 80 + "\r")
        sys.stdout.flush()
    except KeyboardInterrupt:
        print(f"\n\n  {RED}{BOLD}Cancelled!{RESET} Action aborted by user.\n")
        sys.exit(0)

# ── Execute ───────────────────────────────────────────────────────────────────

def execute_action(choice: int):
    if os.name != "nt":
        print(f"\n  {RED}{BOLD}Error:{RESET} {WHITE}This tool only works on Windows.{RESET}\n")
        return

    if choice == 1:    # Shutdown
        subprocess.run(["shutdown", "/s", "/t", "0"], shell=True)
    elif choice == 2:  # Restart
        subprocess.run(["shutdown", "/r", "/t", "0"], shell=True)
    elif choice == 3:  # Sleep
        # SetSuspendState(hibernate, force, disable_wake_events)
        ctypes.windll.PowrProf.SetSuspendState(0, 1, 0)
    elif choice == 4:  # Logout
        subprocess.run(["shutdown", "/l"], shell=True)

# ── Confirmation ──────────────────────────────────────────────────────────────

def confirm(action_name: str, timer_secs: int, battery_target: int = None, auto_confirm: bool = False) -> bool:
    color = {"Shutdown": CYAN, "Restart": GREEN, "Sleep": YELLOW, "Logout": PINK}[action_name]

    print(f"\n  {GRAY}{'─' * 50}{RESET}")
    print(f"  {P2}{BOLD}Summary{RESET}")
    print(f"  {GRAY}{'─' * 50}{RESET}")
    print(f"    {WHITE}Action :{RESET}  {color}{BOLD}{action_name}{RESET}")
    print(f"    {WHITE}Timer  :{RESET}  {format_duration(timer_secs)}")
    if battery_target is not None:
        print(f"    {WHITE}Battery:{RESET}  {YELLOW}Wait until {battery_target}%{RESET}")
    print(f"  {GRAY}{'─' * 50}{RESET}")
    print()

    if auto_confirm:
        return True

    print(f"    {GREEN}{BOLD}[Y]{RESET} {WHITE}Run{RESET}        {RED}{BOLD}[N]{RESET} {WHITE}Go back{RESET}")
    print()

    while True:
        ans = input(f"  {P3}{BOLD}>{RESET} ").strip().lower()
        if ans in ("y", "yes", ""):
            return True
        if ans in ("n", "no", "b", "back"):
            return False
        print(f"  {RED}  Type Y to run or N to go back.{RESET}")

# ── CLI Arguments ─────────────────────────────────────────────────────────────

class CustomArgumentParser(argparse.ArgumentParser):
    def error(self, message):
        if os.name == "nt":
            os.system("")
        print(f"\n  {RED}{BOLD}Error:{RESET} {WHITE}{message}{RESET}")
        
        if "-b/--battery" in message and "expected one argument" in message:
            print(f"  {GRAY}The -b flag requires a battery percentage number.{RESET}")
            print(f"  {GRAY}Example: {RESET}{GREEN}nsleep -s -b 15{RESET} {GRAY}(Shutdown when battery hits 15%){RESET}\n")
        elif "-t/--timer" in message and "expected one argument" in message:
            print(f"  {GRAY}The -t flag requires a time value.{RESET}")
            print(f"  {GRAY}Example: {RESET}{GREEN}nsleep -z -t 45m{RESET} {GRAY}(Sleep in 45 minutes){RESET}\n")
        elif "invalid int value" in message and "-b/--battery" in message:
            print(f"  {GRAY}The -b flag expects a plain number without the % sign.{RESET}")
            print(f"  {GRAY}Example: {RESET}{GREEN}nsleep -s -b 20{RESET}\n")
        elif "--theme" in message and "expected one argument" in message:
            print(f"  {GRAY}The --theme flag requires a theme name.{RESET}")
            print(f"  {GRAY}Available themes: {RESET}{GREEN}lavender, midnight, sunset, forest{RESET}")
            print(f"  {GRAY}Example: {RESET}{GREEN}nsleep --theme sunset{RESET}\n")
        else:
            print(f"  {GRAY}Run nsleep -h for full usage instructions.{RESET}\n")
        sys.exit(2)

def parse_args():
    parser = CustomArgumentParser(description="NythSleep - A stylish Windows 11 power management CLI tool.")

    
    # Action group (mutually exclusive)
    action_group = parser.add_mutually_exclusive_group()
    action_group.add_argument("-s", "--shutdown", action="store_true", help="Power off your machine completely")
    action_group.add_argument("-r", "--restart", action="store_true", help="Reboot your machine")
    action_group.add_argument("-z", "--sleep", action="store_true", help="Put your machine to sleep")
    action_group.add_argument("-l", "--logout", action="store_true", help="Sign out of the current user session")
    
    # Other flags
    parser.add_argument("-t", "--timer", type=str, help="Set a timer (e.g., '1h 30m')")
    parser.add_argument("-i", "--insomnia", action="store_true", help="Keep Awake mode (prevents sleep)")
    parser.add_argument("-b", "--battery", type=int, help="Trigger action at specific battery percentage")
    parser.add_argument("--theme", type=str, choices=["lavender", "midnight", "sunset", "forest"], help="Select UI theme")
    
    return parser.parse_args()

# ── Main ──────────────────────────────────────────────────────────────────────

def main():
    # Enable ANSI on Windows
    if os.name == "nt":
        os.system("")  # enables VT100 escape sequences in cmd/powershell
        
    args = parse_args()
    latest_version = check_for_update()

    # Feature 2: Bypass Mode
    choice = 0
    if args.shutdown: choice = 1
    elif args.restart: choice = 2
    elif args.sleep: choice = 3
    elif args.logout: choice = 4

    if args.theme:
        set_theme(args.theme)

    # Validate that modifiers (-t, -b) are paired with an action or insomnia mode
    if choice == 0 and not args.insomnia:
        if args.timer or args.battery is not None:
            print(f"\n  {RED}{BOLD}Error:{RESET} {WHITE}Modifiers like -t/--timer and -b/--battery require an action flag.{RESET}")
            print(f"  {GRAY}Please specify an action: -s (Shutdown), -r (Restart), -z (Sleep), or -l (Logout).{RESET}")
            print(f"  {GRAY}Example: {RESET}{GREEN}nsleep -s -b 15{RESET}\n")
            return

    if args.insomnia:
        # Validate timer if provided
        timer_secs = 0
        if args.timer:
            timer_secs = parse_timer(args.timer)
            if timer_secs < 0:
                print(f"\n  {RED}{BOLD}Error:{RESET} {WHITE}Invalid timer format '{args.timer}'.{RESET}")
                print(f"  {GRAY}Please use units: {RESET}{CYAN}h{RESET}{GRAY} (hours), {RESET}{CYAN}m{RESET}{GRAY} (minutes), or {RESET}{CYAN}s{RESET}{GRAY} (seconds).{RESET}")
                print(f"  {GRAY}Example: {RESET}{GREEN}nsleep -i -t 30m{RESET}\n")
                return

        set_insomnia_mode(True)
        if choice == 0:
            run_insomnia_loop(latest_version, timer_secs)

    try:
        # If an action flag was passed, bypass the menu
        if choice != 0:
            action_name = ACTIONS[choice - 1][0]
            timer_secs = parse_timer(args.timer) if args.timer else 0
            
            if timer_secs < 0:
                print(f"  {RED}Invalid timer format. Exiting.{RESET}")
                return
                
            clear_screen()
            
            if confirm(action_name, timer_secs, args.battery, auto_confirm=True):
                if args.battery:
                    wait_for_battery(args.battery, action_name)
                countdown(timer_secs, action_name)
                print(f"\n  {P4}{BOLD}Executing {action_name}...{RESET}\n")
                execute_action(choice)
            return

        # Interactive Mode (Default)
        while True:
            clear_screen()
            print_banner(latest_version)
            show_menu()

            choice = get_choice()
            if choice == 0:
                print(f"\n  {P2}Goodbye! Sweet dreams.{RESET}\n")
                break

            action_name = ACTIONS[choice - 1][0]
            timer_secs = get_timer()

            if confirm(action_name, timer_secs, args.battery):
                if args.battery:
                    wait_for_battery(args.battery, action_name)
                countdown(timer_secs, action_name)
                print(f"\n  {P4}{BOLD}Executing {action_name}...{RESET}\n")
                execute_action(choice)
                break
            # else loop back to menu
    except KeyboardInterrupt:
        print(f"\n\n  {P2}Goodbye! Session ended.{RESET}\n")
        sys.exit(0)

if __name__ == "__main__":
    main()
