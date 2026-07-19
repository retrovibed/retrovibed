import 'ddisc.discovery.pb.dart';

// orders candidates by policy_rank ascending (lower is better) - mirrors
// the server's ddisc.Compare (Go), which already folds health/bytes into
// the ranking. id is the final tiebreaker so callers relying on a total
// order (e.g. SplayTreeSet) don't collapse two distinct candidates that
// tie on policy_rank.
int compare(Discovery a, Discovery b) {
  if (a.policyRank != b.policyRank) return a.policyRank.compareTo(b.policyRank);
  return a.id.compareTo(b.id);
}
