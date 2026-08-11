const fs = require('fs');
const path = require('path');

function stripComments(code) {
    let out = '';
    let i = 0;
    let inString = false;
    let stringChar = '';
    let inSingleComment = false;
    let inMultiComment = false;
    let inHtmlComment = false;

    while (i < code.length) {
        let c = code[i];
        let next = code[i + 1];
        let next2 = code[i + 2];
        let next3 = code[i + 3];

        if (inSingleComment) {
            if (c === '\n') {
                inSingleComment = false;
                out += c;
            }
            i++;
            continue;
        }

        if (inMultiComment) {
            if (c === '*' && next === '/') {
                inMultiComment = false;
                i += 2;
            } else {
                if (c === '\n') out += c; // Preserve line numbers slightly better
                i++;
            }
            continue;
        }

        if (inHtmlComment) {
            if (c === '-' && next === '-' && next2 === '>') {
                inHtmlComment = false;
                i += 3;
            } else {
                if (c === '\n') out += c;
                i++;
            }
            continue;
        }

        if (inString) {
            if (c === '\\') {
                out += c;
                if (next !== undefined) {
                    out += next;
                    i += 2;
                } else {
                    i++;
                }
                continue;
            }
            if (c === stringChar) {
                inString = false;
            }
            out += c;
            i++;
            continue;
        }

        // String start
        if (c === '"' || c === "'" || c === '`') {
            inString = true;
            stringChar = c;
            out += c;
            i++;
            continue;
        }

        // Check for comments
        if (c === '/' && next === '/') {
            inSingleComment = true;
            i += 2;
            continue;
        }

        if (c === '/' && next === '*') {
            inMultiComment = true;
            i += 2;
            continue;
        }

        if (c === '<' && next === '!' && next2 === '-' && next3 === '-') {
            inHtmlComment = true;
            i += 4;
            continue;
        }

        out += c;
        i++;
    }
    return out;
}

function processDirectory(dir) {
    const files = fs.readdirSync(dir);
    for (const file of files) {
        const fullPath = path.join(dir, file);
        const stat = fs.statSync(fullPath);
        
        if (stat.isDirectory()) {
            if (file !== 'node_modules' && file !== '.git' && file !== 'dist') {
                processDirectory(fullPath);
            }
        } else {
            const ext = path.extname(file);
            if (['.js', '.ts', '.vue', '.html', '.css', '.scss', '.go'].includes(ext)) {
                if (file === 'remove_frontend_comments.js') continue;
                try {
                    const content = fs.readFileSync(fullPath, 'utf8');
                    const stripped = stripComments(content);
                    if (content !== stripped) {
                        fs.writeFileSync(fullPath, stripped, 'utf8');
                        console.log(`Stripped comments from: ${fullPath}`);
                    }
                } catch (e) {
                    console.error(`Error processing ${fullPath}:`, e);
                }
            }
        }
    }
}

console.log("Starting frontend comment removal...");
processDirectory('.');
console.log("Finished frontend comment removal.");
