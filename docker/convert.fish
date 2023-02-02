#!/usr/bin/env fish

function ee; echo $argv; eval $argv; end

set -gx UID (id -u)
set -gx GID (id -g)
set -gx WORK_DIR $argv[1]
set args  $argv[2..-1]

# echo "            Args: $args"
# echo "  Work Directory: $work_dir"

ee "docker-compose run imagemagick $args"