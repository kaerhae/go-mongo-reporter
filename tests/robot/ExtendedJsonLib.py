from robot.api.deco import keyword
class ExtendedJsonLib:
    @keyword('Filter By Value')
    def FilterByValue(self, json_data, key, value):
        result = [x for x in json_data if x[key] == value]
        return result
    
    @keyword('Check List Length')
    def CheckLength(self, json_data):
        return len(json_data)