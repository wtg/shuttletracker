
# coding: utf-8

# In[ ]:


import requests

url_history = 'https://shuttles.rpi.edu/history'
url_updates = 'https://shuttles.rpi.edu/updates'

r = requests.get(url_history)
j_history = r.json()

r = requests.get(url_updates)
j_updates = r.json()


def mean(L):
    return sum(L)/len(L)


def compute_ave(j):
    speed = []
    for vehicle in j:
        history = [i['speed'] for i in vehicle]
        speed.append( mean(history))
    return mean(speed)


def compute_curr(j):
    if len(j) == 0: ## empty
        return 0
    curr_speed = [i['speed'] for  i in j]
    return mean(curr_speed)


v_ave = compute_ave(j_history)
v_curr = compute_curr(j_updates)

N = len(j_history)

V = ((N-1)*v_ave + v_curr)/N

print ('N = %d, V_ave = %.2f, V_curr = %.2f, V = %.2f' %(N, v_ave, v_curr, V) )

