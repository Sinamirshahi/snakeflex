#!/usr/bin/env python3
"""
Demo Script for Python Web Terminal
Shows various output types including errors
"""

import time
import sys
import random

def main():
    print("🚀 Python Web Terminal Demo Script")
    print("=" * 40)
    
    # Normal output
    print("✅ This is normal stdout output")
    time.sleep(1)
    
    # Error output
    print("⚠️  This is stderr output", file=sys.stderr)
    time.sleep(1)
    
    # Progress simulation
    print("\n📊 Simulating some work...")
    for i in range(5):
        print(f"   Processing step {i+1}/5...")
        time.sleep(0.8)
    
    # Random number generation
    print(f"\n🎲 Random number: {random.randint(1, 100)}")
    
    # File operations demo
    try:
        print("\n📁 Attempting to read a file...")
        with open("nonexistent.txt", "r") as f:
            content = f.read()
    except FileNotFoundError:
        print("❌ File not found (this is expected!)", file=sys.stderr)
    
    # Math operations
    print(f"\n🧮 Some calculations:")
    print(f"   • 2 + 2 = {2 + 2}")
    print(f"   • 10 ** 3 = {10 ** 3}")
    print(f"   • π ≈ {3.14159}")
    
    # Success message
    print(f"\n🎉 Demo completed successfully!")
    print("   Check out both stdout (white) and stderr (red) messages!")

if __name__ == "__main__":
    main()