import os
import shutil
import json
import matplotlib.pyplot as plt
from matplotlib.ticker import PercentFormatter
from matplotlib import ticker
from matplotlib.pyplot import MultipleLocator
import pylab
from datetime import datetime
import click


plt.rcParams["figure.figsize"] = (6, 4.5)
plt.rcParams.update({'font.size': 22})

cur_dir_list = []


all_test_time_list = []
all_unique_bug_num_list = []

all_buggy_test_case = []


def check_is_timeout(cur_dir, test_name_str):
	cur_stdout = os.path.join(cur_dir, "exec", test_name_str, "stdout")
	cur_stdout = open(cur_stdout, 'r')

	for line in cur_stdout.readlines():
		if "panic: test timed out after 30s" in line:
			cur_stdout.close()
			return True

	cur_stdout.close()
	return False

@click.command()
@click.option("--with-feedback-path", required = True)
@click.option("--no-feedback-path", required = True)
@click.option("--no-oracle-path", required = True)
@click.option("--no-mutation-path", required = True)

def main(with_feedback_path, no_feedback_path, no_oracle_path, no_mutation_path):

	cur_dir_list.append(with_feedback_path)
	cur_dir_list.append(no_feedback_path)
	cur_dir_list.append(no_oracle_path)
	cur_dir_list.append(no_mutation_path)


	for cur_dir in cur_dir_list:
		cur_fuzzer_log = open(os.path.join(cur_dir, "fuzzer.log"))

		start_time = None
		uniq_bug_num = 0

		test_time_list = []
		unique_bug_num_list = []

		buggy_test_case = []

		total_lines = cur_fuzzer_log.readlines()
		total_lines_num = len(total_lines)

		cur_test_case = ""
		cur_test_case_sim = ""

		prev_duration = -0.2

		for idx, cur_line in enumerate(total_lines):

			if cur_line[:2] != "20":
				continue

			if "received " in cur_line:
				cur_test_case = cur_line.split("received ")[-1]
				cur_test_case = cur_test_case.replace("\n", "")
				cur_test_case_sim = cur_test_case.split("-")[-2]

			time_str = cur_line.split(" ")
			time_str = time_str[0] + " " + time_str[1]
			cur_time = datetime.strptime(time_str, '%Y/%m/%d %H:%M:%S')

			if start_time == None:
				start_time = cur_time
				continue

			cur_duration = (cur_time - start_time).total_seconds()
			# To hours
			cur_duration = cur_duration / 3600

			if idx == total_lines_num - 1:
				test_time_list.append(cur_duration)
				if len(unique_bug_num_list) > 0:
					unique_bug_num_list.append(unique_bug_num_list[-1])
				else:
					unique_bug_num_list.append(0)


			# Restrict time frame
			if cur_duration > 12.0:
				test_time_list.append(cur_duration)
				if len(unique_bug_num_list) > 0:
					unique_bug_num_list.append(unique_bug_num_list[-1])
				else:
					unique_bug_num_list.append(0)
				break

			if "unique bug(s)" in cur_line and "found" in cur_line:
				# Ignore timeout cases
				# if check_is_timeout(cur_dir, cur_test_case):
				# 	print("For test case: %s : %s, timeout. " % (cur_dir, cur_test_case))
				# 	continue

				# test_time_list.append(cur_duration)
				# unique_bug_num_list.append(uniq_bug_num)
				cur_unique_bug_num_str = cur_line.split(" unique bug(s)")[0]
				cur_unique_bug_num_str = cur_unique_bug_num_str.split("found ")[1]
				uniq_bug_num += int(cur_unique_bug_num_str)
				# test_time_list.append(cur_duration)
				# unique_bug_num_list.append(uniq_bug_num)

				continue

			if cur_duration - prev_duration >= 0.03:
				test_time_list.append(cur_duration)
				unique_bug_num_list.append(uniq_bug_num)
				prev_duration = cur_duration



		all_test_time_list.append(test_time_list)
		all_unique_bug_num_list.append(unique_bug_num_list)

	print("Log: %d" % (all_test_time_list[0][-1]))
	print("Log: The last unique bug number is: %d. " % (all_unique_bug_num_list[0][-1]))
	print("Log: The last unique bug number is: %d. " % (all_unique_bug_num_list[1][-1]))

	x_major_locator=MultipleLocator(0.5)
	# y_major_locator=MultipleLocator(10)

	plt.figure()
	ax = plt.subplot()

	ax.plot(all_test_time_list[0], all_unique_bug_num_list[0], linestyle = 'dashed', marker = 'p', markevery=6, linewidth=2.0, markersize=7)
	ax.plot(all_test_time_list[1], all_unique_bug_num_list[1], linestyle = (0, (3, 1, 1, 1, 1, 1)), marker = 'o', markevery=6, linewidth=2.0, markersize=7)
	ax.plot(all_test_time_list[2], all_unique_bug_num_list[2], linestyle = 'dashdot', marker = 's', markevery=6, linewidth=2.0, markersize=7)
	ax.plot(all_test_time_list[3], all_unique_bug_num_list[3], linestyle = 'dotted', marker = '*', markevery=6, linewidth=2.0, markersize=7)
	#ax.plot(all_test_time_list[4], all_unique_bug_num_list[4], linestyle = 'dashed', marker = '^', markevery=4, linewidth=2.0, markersize=9)
	#ax.plot(all_test_time_list[5], all_unique_bug_num_list[5], linestyle = (0, (3, 1, 1, 1, 1, 1)), marker = 'v', markevery=4, linewidth=2.0, markersize=9)
	#ax.plot(all_test_time_list[6], all_unique_bug_num_list[6], linestyle = 'dotted', marker = 'X', markevery=4, linewidth=2.0, markersize=9)
	#ax.plot(all_test_time_list[7], all_unique_bug_num_list[7], linestyle = 'dashed', marker = 'D', markevery=4, linewidth=2.0, markersize=9)

	plt.title("Contribution of GFuzz components", fontsize=20)
	# plt.title("GFuzz rand stage only", fontsize=20)
	plt.xlabel("Time (h)", fontsize=20)
	plt.ylabel("Num of Unique Bugs", fontsize=20)
	plt.xticks(fontsize=20)
	plt.yticks(fontsize=20)
	leg = plt.legend(['GFuzz', 'no_feedbacks', 'no_oracle', 'no_mutations'], fontsize=14, handlelength=3)
	plt.xlim([0,3])
	plt.ylim([0, 20])
	ax.xaxis.set_major_locator(x_major_locator)
	# ax.yaxis.set_major_locator(y_major_locator)

	plt.grid()

	plt.tight_layout()
	plt.savefig('./bug_vs_time.png', dpi = 200)

	# fig = plt.figure()
	# ax = fig.add_subplot(111)
	# x = [0, 1]
	# y = [0, 1]
	# ax.plot(x, y, linestyle = 'dashed', marker = 'p', markevery=2, linewidth=2.0, markersize=9)
	# # ax.plot(x, y, linestyle = (0, (3, 1, 1, 1, 1, 1)), marker = 'o', markevery=2, linewidth=2.0, markersize=9)
	# #ax.plot(x, y, linestyle = 'dashdot', marker = 's', markevery=2, linewidth=2.0, markersize=9, label=legends[2])
	# #ax.plot(x, y, linestyle = 'dotted', marker = '*', markevery=2, linewidth=2.0, markersize=9, label=legends[3])
	# #ax.plot(x, y, linestyle = 'dashed', marker = '^', markevery=4, linewidth=2.0, markersize=9, label=legends[4])
	# #ax.plot(x, y, linestyle = (0, (3, 1, 1, 1, 1, 1)), marker = 'v', markevery=4, linewidth=2.0, markersize=9, label=legends[5])
	# #ax.plot(x, y, linestyle = 'dotted', marker = 'X', markevery=4, linewidth=2.0, markersize=9, label=legends[6])
	# #ax.plot(x, y, linestyle = 'dashed', marker = 'D', markevery=4, linewidth=2.0, markersize=9, label=legends[7])
	# # save it *without* adding a legend
	# # fig.savefig('image.png')

	# figsize = (6, 4.5)
	# fig_leg = plt.figure(figsize=figsize)
	# ax_leg = fig_leg.add_subplot(111)
	# # add the legend from the previous axes
	# ax_leg.legend(*ax.get_legend_handles_labels(), loc='center', ncol = 2, handlelength=3, fontsize = 14)
	# # hide the axes frame and the x/y labels
	# ax_leg.axis('off')
	# fig_leg.savefig('./legend.png')

if __name__ == '__main__':
	main()
