function character_said(text)
  if text == "shop" then
    shop("WOBGI_TAVERN")
  end
end

function interact()
  say("What is it you'd like to learn more about?")
  convo_select("weather", "Can you tell me about the weather here?")
  convo_select("whereami", "What is this place?")
end

function conversation_select(option_id)
  if option_id == "weather" then
    say("The weather is very odd.. to say the least.")
  elseif option_id == "whereami" then
    say("I wish I knew.")
  end
end
