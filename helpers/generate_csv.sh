#!/bin/bash

OUTPUT_FILE=${2:-"random_data.csv"}
MAX_COL=${1:-10}
NUM_ROWS=1000000

awk -v cols="$MAX_COL" -v rows="$NUM_ROWS" '
    BEGIN {
        srand();
        OFS = ","
        
        # 1. Gen antet: indice1,indice2....
        for (i = 1; i <= cols; i++) {
            printf "indice%d%s", i, (i == cols ? ORS : OFS)
        }

        # Doar litere mari și cifre (36 de caractere în total)
        split("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", chars, "")

        # 2. Generare date
        for (r = 1; r <= rows; r++) {       
            for (c = 1; c <= cols; c++) {
                
                # Gen 17 caractere
                val = ""
                for (len = 1; len <= 17; len++) {
                    # 36 este lungimea vectorului chars
                    val = val chars[int(rand() * 36) + 1]
                }
                
                printf "%s%s", val, (c == cols ? ORS : OFS)
            }
        }
    }
' > "$OUTPUT_FILE"