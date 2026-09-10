export function mapTeamSettings(teamInfo, members, inviteCode) {
  const leader = members.find((member) => member.id === teamInfo.captain_id);
  return {
    name: teamInfo.name,
    picture: teamInfo.picture,
    description: teamInfo.description,
    captainId: teamInfo.captain_id,
    inviteCode,
    leader: {
      name: leader?.name || '',
      picture: leader?.picture || '',
      email: leader?.email || '',
    },
    members: members
      .filter((member) => member.id !== teamInfo.captain_id)
      .map((member) => ({
        id: member.id,
        name: member.name,
        picture: member.picture,
        email: member.email,
      })),
  };
}

export function teamUpdatePayload(team) {
  return { name: team.name, description: team.description };
}

// A refresh must not overwrite fields changed after it started.
export function unchangedTeamFields(data, startedVersions, currentVersions) {
  return Object.fromEntries(Object.entries(data).filter(([key]) => startedVersions[key] === currentVersions[key]));
}
