#include <vector>
using namespace std;

class Solution {
public:
    vector<int> findDisappearedNumbers(vector<int>& nums) {
        vector<int> result_nums;
        vector<bool> exist(nums.size()+1, false);

        for (int num : nums){
            exist[num] = true;
        }

        for (int i = 1; i < exist.size(); i++){
            if(!exist[i]){
                result_nums.push_back(i);
            }
        }
        return result_nums;
    }
};