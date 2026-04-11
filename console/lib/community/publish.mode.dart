import 'package:retrovibed/billing/plan.summary.dart';
import 'package:retrovibed/community/community.pb.dart';

PublishMode maxPublishMode(String planId) {
  final List<PlanSummary> syndicated = [family(), premium(), founder()];
  final List<PlanSummary> listable = [personal4(), personal()];
  if (syndicated.any((p) => p.id == planId)) return PublishMode.SYNDICATED;
  if (listable.any((p) => p.id == planId)) return PublishMode.LISTED;
  return PublishMode.UNLISTED;
}
