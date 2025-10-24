#!/usr/bin/env python3
"""Test script to verify project structure without dependencies."""
import os
import sys

def test_project_structure():
    """Test that all required files and directories exist."""
    
    required_files = [
        'README.md',
        'requirements.txt',
        '.gitignore',
        '.env.example',
        'config.py',
        'main.py',
        'meme_generator.py',
        'discord_bot.py',
        'telegram_bot.py',
        'example.py',
        'QUICKSTART.md'
    ]
    
    required_dirs = [
        'meme_templates',
        'output'
    ]
    
    print("Testing Meme Bot Project Structure")
    print("=" * 50)
    
    # Check files
    print("\nChecking required files:")
    all_files_exist = True
    for file in required_files:
        exists = os.path.exists(file)
        status = "✅" if exists else "❌"
        print(f"  {status} {file}")
        if not exists:
            all_files_exist = False
    
    # Check directories
    print("\nChecking required directories:")
    all_dirs_exist = True
    for directory in required_dirs:
        exists = os.path.isdir(directory)
        status = "✅" if exists else "❌"
        print(f"  {status} {directory}/")
        if not exists:
            all_dirs_exist = False
    
    # Check Python syntax
    print("\nChecking Python files syntax:")
    python_files = [f for f in required_files if f.endswith('.py')]
    syntax_ok = True
    
    for py_file in python_files:
        try:
            with open(py_file, 'r') as f:
                compile(f.read(), py_file, 'exec')
            print(f"  ✅ {py_file}")
        except SyntaxError as e:
            print(f"  ❌ {py_file}: {e}")
            syntax_ok = False
    
    # Summary
    print("\n" + "=" * 50)
    if all_files_exist and all_dirs_exist and syntax_ok:
        print("✅ Project structure is complete and valid!")
        print("\nNext steps:")
        print("1. Install dependencies: pip install -r requirements.txt")
        print("2. Configure bot token: cp .env.example .env (then edit)")
        print("3. Add meme templates to meme_templates/")
        print("4. Run bot: python main.py discord OR python main.py telegram")
        return 0
    else:
        print("❌ Project structure has issues!")
        if not all_files_exist:
            print("  - Some required files are missing")
        if not all_dirs_exist:
            print("  - Some required directories are missing")
        if not syntax_ok:
            print("  - Some Python files have syntax errors")
        return 1


if __name__ == '__main__':
    sys.exit(test_project_structure())
