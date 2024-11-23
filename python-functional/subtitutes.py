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
            if line_number >= len(all_lines):
                raise Exception("Invalid line number")
            new_list = all_lines[:]
            if line_number == len(all_lines):
                new_list.append("\n")
            else:
                original = all_lines[line_number]
                new_list[line_number] = f"{original}\n"
            return "\n".join(new_list)
        case EditType.SUBSTITUTE:
            line_number = edit["line_number"]
            insert_text = edit["insert_text"]
            start = edit["start"]
            end = edit["end"]
            all_lines = document.splitlines()
            
            if line_number >= len(all_lines):
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
        
        case EditType.DELETE:
            line_number = edit["line_number"]
            start = edit["start"]
            end = edit["end"]
            all_lines = document.splitlines()        
            if line_number > len(all_lines):
                raise Exception("Invalid line number")
            toreplace = all_lines[line_number]
            if start > len(toreplace):
                raise Exception("Invalid start index")
            if (end < start or end >len(toreplace)):
                raise Exception("Invalid end index")
            prev = toreplace[:start]
            next = toreplace[end:]
            all_lines[line_number]=prev + next
            return "\n".join(all_lines)

        case EditType.INSERT:
            line_number = edit["line_number"]
            insert_text = edit["insert_text"]
            start = edit["start"]
            all_lines = document.splitlines()
            print(all_lines)
            print(line_number)
            if line_number > len(all_lines):
                raise Exception("Invalid line number")
            nws_list = all_lines[:]
            
            if line_number == len(all_lines):
                # If inserting at the end, just append the new text
                nws_list.append(insert_text)
            else:
                # Otherwise modify existing line
                toreplace = nws_list[line_number]
                if start > len(toreplace):
                    raise Exception("Invalid start index")
                new_line = toreplace[:start] + insert_text + toreplace[start:]
                nws_list[line_number] = new_line
            
            return "\n".join(nws_list)
        case default:
            raise Exception("Unknown edit type")

