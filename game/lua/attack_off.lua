if not ma then
    for i = 1,72 do
        if IsAttackAction(i) then
             ma = i;
        end;
    end;
end; 
if ma then
    if IsCurrentAction(ma) then
        AttackTarget("target");
    end;
end;