if not ma then
    for i = 1,72 do
        if IsAttackAction(i) then
             ma = i;
        end;
    end;
end; 
if ma then
    if not IsCurrentAction(ma) then
        UseAction(ma);
    end; 
    else AttackTarget("target");
end;