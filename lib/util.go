package lib

import "github.com/diamondburned/arikawa/v3/discord"

func GuildPermissionsOf(guild discord.Guild, member discord.Member) discord.Permissions {
  if guild.OwnerID == member.User.ID {
    return discord.PermissionAll
  }

  var perm discord.Permissions

  for _, role := range guild.Roles {
    if role.ID == discord.RoleID(guild.ID) {
      perm |= role.Permissions
      break
    }
  }

  if perm.Has(discord.PermissionAdministrator) {
    return discord.PermissionAll
  }

  for _, role := range guild.Roles {
    for _, id := range member.RoleIDs {
      if id == role.ID {
        perm |= role.Permissions
      }
    }
  }

  if perm.Has(discord.PermissionAdministrator) {
    return discord.PermissionAll
  }

  return perm
}

func ChannelPermissionsOf(guild discord.Guild, channel discord.Channel, member discord.Member) discord.Permissions {
  perm := GuildPermissionsOf(guild, member)

  if perm.Has(discord.PermissionAdministrator) {
    return discord.PermissionAll
  }

  for _, overwrite := range channel.Overwrites {
    if discord.GuildID(overwrite.ID) == guild.ID {
      perm &= ^overwrite.Deny
      perm |= overwrite.Allow
      break
    }
  }

  var deny, allow discord.Permissions

  for _, overwrite := range channel.Overwrites {
    for _, id := range member.RoleIDs {
      if id == discord.RoleID(overwrite.ID) && overwrite.Type == discord.OverwriteRole {
        deny |= overwrite.Deny
        allow |= overwrite.Allow
      }
    }
  }

  perm &= ^deny
  perm |= allow

  for _, overwrite := range channel.Overwrites {
    if discord.UserID(overwrite.ID) == member.User.ID && overwrite.Type == discord.OverwriteMember {
      perm &= ^overwrite.Deny
      perm |= overwrite.Allow
    }
  }

  if perm.Has(discord.PermissionAdministrator) {
    return discord.PermissionAll
  }

  return perm
}
