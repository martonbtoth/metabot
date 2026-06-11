local i = 1
all_spells = ""
while true do
   local spellName, spellRank = GetSpellName(i, BOOKTYPE_SPELL)
   if not spellName then
      do break end
   end
   
   -- use spellName and spellRank here
   all_spells = all_spells .. spellName .. '(' .. spellRank .. ')' .. '\n'
   
   i = i + 1
end

{0} = all_spells