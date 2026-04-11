import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/billing/plan.summary.dart';
import 'package:retrovibed/community/community.pb.dart';
import 'package:retrovibed/community/publish.mode.dart';

void main() {
  group('maxPublishMode', () {
    test('free plan returns UNLISTED', () {
      expect(maxPublishMode(free().id), equals(PublishMode.UNLISTED));
    });

    test('personal plan returns LISTED', () {
      expect(maxPublishMode(personal().id), equals(PublishMode.LISTED));
    });

    test('personal4 plan returns LISTED', () {
      expect(maxPublishMode(personal4().id), equals(PublishMode.LISTED));
    });

    test('family plan returns SYNDICATED', () {
      expect(maxPublishMode(family().id), equals(PublishMode.SYNDICATED));
    });

    test('premium plan returns SYNDICATED', () {
      expect(maxPublishMode(premium().id), equals(PublishMode.SYNDICATED));
    });

    test('founder plan returns SYNDICATED', () {
      expect(maxPublishMode(founder().id), equals(PublishMode.SYNDICATED));
    });

    test('unknown plan returns UNLISTED', () {
      expect(maxPublishMode('unknown-plan-id'), equals(PublishMode.UNLISTED));
    });

    test('empty plan id returns UNLISTED', () {
      expect(maxPublishMode(''), equals(PublishMode.UNLISTED));
    });
  });
}
