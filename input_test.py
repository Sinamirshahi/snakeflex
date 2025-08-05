#!/usr/bin/env python3
"""
Simple Input Test Script
Basic test for interactive input functionality
"""
import sys

def main():
    print("Simple Input Test")
    print("================")
    print()  # Add blank line
    
    # Test basic input - FLUSH OUTPUT BEFORE INPUT
    print("Enter your name: ", end="", flush=True)
    name = input()
    print(f"Hello, {name}!")
    print()
    
    # Test number input
    print("Enter a number: ", end="", flush=True)
    try:
        number = int(input())
        print(f"Your number doubled is: {number * 2}")
    except ValueError:
        print("That wasn't a valid number!")
    
    print()
    print("Test completed!")

if __name__ == "__main__":
    main()