from enum import Enum


class EditType(Enum):
    NEWLINE = 1
    SUBSTITUTE = 2
    INSERT = 3
    DELETE = 4


def handle_edit(document, edit_type, edit):
    match edit_type:
        case EditType.NEWLINE:
            line_number = edit["line_number"]
            all_lines = document.splitlines()
            if line_number > len(all_lines):
                raise Exception("Invalid line number")
            original = all_lines[line_number]
            all_lines[line_number] = f"{original}\n"
            return "\n".join(all_lines)
        case EditType.SUBSTITUTE:
            line_number = edit["line_number"]
            insert_text = edit["insert_text"]
            start = edit["start"]
            end = edit["end"]
            all_lines = document.splitlines()
            
            if line_number > len(all_lines):
                raise Exception("Invalid line number")
            toreplace = all_lines[line_number]
            print("line", toreplace)
            if start > len(toreplace):
                raise Exception("Invalid start index")
            if (end < start or end >len(toreplace)):
                raise Exception("Invalid end index")
            new_line = toreplace[:start] + insert_text + toreplace[end:]
            all_lines[line_number] = new_line
            return "\n".join(all_lines)
            
            
    

