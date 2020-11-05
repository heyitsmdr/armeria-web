-- comment
function character_entered(char_name)
  say("Hey there, " .. char_name)
end

function character_said(char_name, text)
  if text == "hey" then
    say("Greetings!")
  elseif text == "title" then
    say("Your title is " .. c_attr(char_name, "title", false))
  elseif text == "set" then
    s = c_set_attr(char_name, "title", "Haha", false)
    say("The function returned " .. s)
  end
end