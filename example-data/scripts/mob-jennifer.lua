-- comment
function character_entered(char_name)
  say("Hey there, " .. char_name)
end

function character_said(char_name, text)
  if text == "hey" then
    say("Greetings!")
  elseif text == "title" then
    say(c_attr(char_name, "title", false))
  end
end