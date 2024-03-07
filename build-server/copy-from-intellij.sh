#!/bin/bash

dir_scratch="/home/kiquetal/.config/JetBrains/IntelliJIdea2023.3/scratches"
out_dir="./scratches"

# Create the output directory if it doesn't exist
mkdir -p "$out_dir"

# Copy all files from the scratch directory to the output directory
cp "$dir_scratch"/*.http "$out_dir"

