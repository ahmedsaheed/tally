#!/usr/bin/env bash
directories=($(ls -d ~/Developer/*/))
random_directory=${directories[$RANDOM % ${#directories[@]}]}

# Echo results from running rh tally
echo "Running tally in $random_directory"
echo $(go run . $random_directory)


# Echo results from running rh kc
echo "Running kc in $random_directory"
echo $(kc $random_directory)


echo "Running benchmarking script in $random_directory"
hyperfine 'kc $random_directory' 'tally $random_directory' -i --export-markdown ./performance.md 


glow ./performance.md
