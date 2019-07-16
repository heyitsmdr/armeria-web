-- comment
function character_entered(char_name)
  say("Hey there, " .. char_name)
end

function character_said(text)
  if text == "hey" then
    say("Greetings!")
  end
end