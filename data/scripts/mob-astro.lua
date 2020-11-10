-- Astro Script
function character_entered()
  say("Welcome to the Armeria Development Environment, " .. invoker_name .. "!")
  say("You are on [b]Test Island[/b]. You can do all of your testing here.")
  say("For more information, say [b]help[/b].")
end

function character_said(text)
  if text == "help" then
    say("What would you like to know more about? You can ask me about [b]test island[/b].")
  elseif text == "title" then
    say("Your title is " .. c_attr(char_name, "title", false))
  elseif text == "shop" then
    shop()
  end
end