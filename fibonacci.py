#!/usr/bin/env python3
"""
Fibonacci Sequence Generator
A demonstration script for the Python Web Terminal
"""

import time
import sys

def fibonacci_generator(n):
    """Generate fibonacci sequence up to n numbers"""
    a, b = 0, 1
    for i in range(n):
        yield a
        a, b = b, a + b

def print_banner():
    """Print a fancy banner"""
    banner = """
╔══════════════════════════════════════╗
║        🐍 FIBONACCI GENERATOR 🐍      ║
║                                      ║
║  Generating beautiful number         ║
║  sequences since 1202!               ║
╚══════════════════════════════════════╝
    """
    print(banner)

def main():
    print_banner()
    print("Starting Fibonacci sequence generation...")
    print("=" * 50)
    
    # Get number of fibonacci numbers to generate
    try:
        n = 15  # Default value
        print(f"Generating first {n} Fibonacci numbers:")
        print("-" * 30)
        
        # Generate and display fibonacci numbers with a slight delay for visual effect
        for i, fib_num in enumerate(fibonacci_generator(n)):
            print(f"F({i:2d}) = {fib_num:>8,}")
            time.sleep(0.3)  # Small delay to see streaming effect
            
        print("-" * 30)
        print("✅ Fibonacci sequence generation complete!")
        
        # Calculate some interesting facts
        fib_list = list(fibonacci_generator(n))
        print(f"\n📊 Statistics:")
        print(f"   • Total numbers generated: {len(fib_list)}")
        print(f"   • Largest number: {max(fib_list):,}")
        print(f"   • Sum of all numbers: {sum(fib_list):,}")
        
        # Golden ratio approximation
        if n > 2:
            golden_ratio = fib_list[-1] / fib_list[-2]
            print(f"   • Golden ratio approximation: {golden_ratio:.6f}")
            print(f"   • Actual golden ratio: {(1 + 5**0.5) / 2:.6f}")
        
        print(f"\n🎉 Mission accomplished! The numbers are beautiful, aren't they?")
        
    except KeyboardInterrupt:
        print("\n\n⚠️  Script interrupted by user")
        sys.exit(1)
    except Exception as e:
        print(f"\n❌ An error occurred: {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()