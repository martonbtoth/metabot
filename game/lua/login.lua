if (WoWAccountSelectDialog and WoWAccountSelectDialog:IsShown()) then
	for i = 0, GetNumGameAccounts() do
		if GetGameAccountInfo(i) == '$accname' then
			WoWAccountSelect_SelectAccount(i)
		end
	end
elseif (AccountLoginUI and AccountLoginUI:IsVisible()) then
	DefaultServerLogin('$username', '$password')
	AccountLoginUI:Hide()
elseif (RealmList and RealmList:IsVisible()) then
	for i = 1, select('#',GetRealmCategories()) do
		for j = 1, GetNumRealms(i) do
			if GetRealmInfo(i, j) == '$realm' then
				RealmList:Hide()
				ChangeRealm(i, j)
			end
		end
	end
elseif (CharacterSelectUI and CharacterSelectUI:IsVisible()) then
	if GetServerName() ~= '$realm' and (not RealmList or not RealmList:IsVisible()) then
		RequestRealmList(1)
	else
		for i = 0,GetNumCharacters() do
			if (GetCharacterInfo(i) == '" + name + "') then
				CharacterSelect_SelectCharacter(i)
				EnterWorld()
			end
		end
	end
end