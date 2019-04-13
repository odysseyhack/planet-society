//
//  TransactionFlowContainerViewController.swift
//  Permission Hub
//
//  Created by Corné on 12/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

struct TransactionFlowStep {
    let viewController = UIViewController()
}

final class TransactionFlowContainerViewController: UIPageViewController {

    // MARK: - Private properties

    private let steps: [TransactionFlowStep]
    private var currentStepIndex = 0

    private let service = NetworkingService()

    // MARK: - Initialization

    init(steps: [TransactionFlowStep]) {
        self.steps = steps

        super.init(
            transitionStyle: .scroll,
            navigationOrientation: .horizontal,
            options: [:])
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Life cycle

    override func viewDidLoad() {
        super.viewDidLoad()

        view.backgroundColor = PHColors.lightGray

        // configure navigation bar
        navigationController?.navigationBar.tintColor = PHColors.greyishBrown
        navigationItem.backBarButtonItem = UIBarButtonItem(title: "", style: .plain, target: nil, action: nil)

        // set initial viewcontroller
        let firstViewControllerOrNil = self.steps.map { $0.viewController }.first
        if let viewController = firstViewControllerOrNil {

            setViewControllers(
                [viewController],
                direction: .forward,
                animated: false)
        }
    }

    override func viewDidAppear(_ animated: Bool) {
        super.viewDidAppear(animated)

        pollForNotifications()
    }

    // MARK: - Networking

    @objc private func pollForNotifications() {

        try! service.getNotifications { result in

            switch result {
            case .success(let transactionOrNil):
                if let transaction = transactionOrNil {
                    self.presentTransactionOverviewViewController(transaction: transaction)
                }
            case .failure(let error):
                print(error)
            }
        }

        // poll endpoint every second
        perform(#selector(pollForNotifications), with: nil, afterDelay: 1)
    }

    // MARK: - Helpers

    private func presentTransactionOverviewViewController(transaction: TransactionNotification) {

        // convert items
        var items = [PHTableViewViewCellType]()
        items.append(.notification(
            type: .notification,
            text: "Verified"))
        items.append(.notification(
            type: .warning,
            text: "Permission warning!"))
        items.append(.description(
            title: transaction.title,
            description: transaction.description))

        transaction.items.forEach {
            items.append(.transactionItem(item: $0))
        }

        let viewController = PHTableViewController(items: items)
        navigationController?.pushViewController(viewController, animated: true)
    }
}
