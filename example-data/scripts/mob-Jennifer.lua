-- comment
function character_entered(char_name)
  mob_say("Hey there, " .. char_name)
end

function character_said(text)
  if text == "hey" then
    mob_say "Greetings!"
  end
end