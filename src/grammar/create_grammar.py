"""
Autogenerates the grammar.c file from grammar.csv file.
"""
from _ast import Dict, Tuple
from dataclasses import dataclass
from typing import List, Optional

@dataclass
class Rule:
    """
    This emulates a rule with an op code and its corresponding chars.
    """
    chars: List[Optional[int]]
    op_code: str

    @staticmethod
    def parse_rule(line: str) -> "Rule":
        """
        Parses a rule given a line.
        :param line: A line that correspondes to the rule in the csv file.
        :return: Returns the rule.
        """
        chars_ = line.split(", ")
        chars_hex: List[Optional[int]] = []
        op_code: str = ""
        for char_ in chars_:
            char_ = char_.strip()
            if char_ != "N" and len(char_) == 1:
                chars_hex.append(int(char_, base=16))
            elif char_ == "N":
                chars_hex.append(None)
            else:
                op_code = char_
        return Rule(chars_hex, op_code)


def generate_rules(file_name: str) -> List[Rule]:
    """
    Generate a list of rules given a csv files.
    :param file_name: File name of the csv.
    :return: Returns the list of parsed rule.
    """
    with open(file_name, "r") as file:
        data = file.read()
    rules: List[Rule] = []
    for line in data.split("\n"):
        if line == "":
            continue
        rules.append(Rule.parse_rule(line))
    return rules


def create_grammar(rules: List[Rule]) -> Tuple(dict, List[str]):
    """
    Create the rules and the op_code list.
    :param rules: Rules as a list of Rule object.
    :return: Dictionary that recursively describes an op_code instruction and a list of op codes.
    """
    op_codes: List[str] = []
    grammar: Dict[Optional[int], Dict[Optional[int],   Dict[Optional[int], Dict[Optional[int], str]]]] = {}
    for rule in rules:
        op_codes.append(rule.op_code)
        if grammar.get(rule.chars[0]) is None:
            grammar[rule.chars[0]] = {rule.chars[1]: {rule.chars[2]: {}}}
        elif grammar[rule.chars[0]].get(rule.chars[1]) is None:
            grammar[rule.chars[0]][rule.chars[1]] = {rule.chars[2]: {}}
        elif grammar[rule.chars[0]][rule.chars[1]].get(rule.chars[2]) is None:
            grammar[rule.chars[0]][rule.chars[1]][rule.chars[2]] = {}
        grammar[rule.chars[0]][rule.chars[1]][rule.chars[2]][rule.chars[3]] = rule.op_code
    return grammar, op_codes


def create_switch_case(grammar: dict, depth: int) -> str:
    """
    Create a switch statement given a recursive grammar rule.
    :param grammar: Recursively described grammar as a dictionary.
    :param depth: Depth of the current switch.
    :return: String representation of a switch statement.
    """
    switch_statement = f"{'    ' * (depth + 1)}switch (instruction->characters[{depth}])" + " {"
    for lexeme in grammar:
        if lexeme == None and depth == 3:
            switch_statement += f"\n{'    ' * (depth + 2)}default: return {grammar[lexeme]}; break;\n"
        elif lexeme != None and depth == 3:
            switch_statement += f"\n{'    ' * (depth + 2)}case {lexeme}: return {grammar[lexeme]}; break;\n"
        elif lexeme == None:
            switch_statement += f"\n{'    ' * (depth + 2)}default:\n{create_switch_case(grammar[lexeme], depth + 1)}\n{'    ' * (depth + 1)}break;\n"
        else:
            switch_statement += f"\n{'    ' * (depth + 2)}case {lexeme}:\n{create_switch_case(grammar[lexeme], depth + 1)}\n{'    ' * (depth + 1)}break;\n"
    switch_statement += f"{'    ' * (depth + 1)}" + "}"
    return switch_statement


def create_enum(op_codes: List[str]) -> str:
    """
    Create the op_code type object as an enum.
    :param op_codes: List of op codes as strings.
    :return: String representation of a typdef enum statement.
    """
    enum = "typedef enum {\n    " + ',\n    '.join(op_codes) + "\n} chip8_op_code;"
    return enum


def create_parse_func(grammar: dict) -> str:
    """
    Create the string representation of the parse function given grammar.
    :param grammar: Grammar rules as a dictionary recursively describing the rules.
    :return: The string representation of the parse function.
    """
    func = "chip8_op_code parse_func(chip8_instruction* instruction) {\n" + create_switch_case(grammar, 0) + "\n}"
    return func


def create_file(grammar: dict, op_codes: List[str]) -> None:
    """
    Create the file given grammar rules and list of op codes.
    :param grammar: Grammar rules as a dictionary recursively describing the rules.
    :param op_codes: The string representation of the parse function.
    :return: None, but create the file grammar.c
    """
    with open("grammar.c", "w") as file:
        file.write(create_enum(op_codes) + "\n\n" + create_parse_func(grammar))


if __name__ == '__main__':
    grammar, op_codes = create_grammar(generate_rules("grammar.csv"))
    create_file(grammar, op_codes)

